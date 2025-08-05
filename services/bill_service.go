package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/JZ23-2/splitbill-backend/database"
	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/JZ23-2/splitbill-backend/handlers"
	"github.com/JZ23-2/splitbill-backend/models"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/google/uuid"
)

func CreateBillWithoutParticipant(req dtos.CreateBillWithoutParticipantRequest) (*dtos.CreateBillWithoutParticipantResponse, error) {

	parsedDate, err := time.Parse("2006-01-02", req.BillDate)
	if err != nil {
		return nil, fmt.Errorf("invalid billDate: %w", err)
	}

	billTaxInt := utils.FormatUSDtoInt(req.Tax + req.Service)
	bill := models.Bill{
		StoreName: req.StoreName,
		Tax:       billTaxInt,
		CreatorID: req.CreatorID,
		CreatedAt: time.Now(),
		BillDate:  parsedDate,
	}

	if err := database.DB.Create(&bill).Error; err != nil {
		return nil, errors.New("failed to create bill: " + err.Error())
	}

	var itemResponses []dtos.CreateBillWithoutParticipantItemResponse

	for _, item := range req.Items {
		itemID := uuid.NewString()
		priceInt := utils.FormatUSDtoInt(item.Price)
		newItem := models.Item{
			ItemID:   itemID,
			BillID:   bill.BillID,
			Name:     item.Name,
			Price:    priceInt,
			Quantity: item.Quantity,
		}

		if err := database.DB.Create(&newItem).Error; err != nil {
			return nil, errors.New("failed to create item: " + err.Error())
		}

		itemResponses = append(itemResponses, dtos.CreateBillWithoutParticipantItemResponse{
			ItemID:    itemID,
			Name:      item.Name,
			Quantity:  item.Quantity,
			UnitPrice: float64(priceInt) / 100.0,
		})
	}

	resp := &dtos.CreateBillWithoutParticipantResponse{
		BillID:    bill.BillID,
		StoreName: bill.StoreName,
		BillDate:  bill.BillDate,
		Tax:       float64(billTaxInt) / 100.0,
		CreatedAt: bill.CreatedAt.Format(time.RFC3339),
		CreatorID: bill.CreatorID,
		Items:     itemResponses,
	}

	return resp, nil
}

func GetBillsByCreator(creatorID string, billID string) ([]dtos.GetBillByCreatorResponse, error) {
	if creatorID == "" {
		return nil, errors.New("creatorId is required")
	}

	var bills []models.Bill
	query := database.DB.Preload("Items").Where("creator_id = ?", creatorID)
	if billID != "" {
		query = query.Where("bill_id = ?", billID)
	}
	if err := query.Find(&bills).Error; err != nil {
		return nil, errors.New("failed to fetch bills: " + err.Error())
	}

	var response []dtos.GetBillByCreatorResponse

	for _, bill := range bills {
		var itemResponses []dtos.GetBillByCreatorItemResponse
		participantMap := make(map[string]dtos.GetBillByCreatorParticipantResponse)

		subTotal := 0
		for _, item := range bill.Items {
			subTotal += item.Price
		}
		if subTotal == 0 {
			continue
		}

		for _, item := range bill.Items {
			var participants []models.Participant
			if err := database.DB.Where("item_id = ?", item.ItemID).Find(&participants).Error; err != nil {
				return nil, errors.New("failed to fetch participants: " + err.Error())
			}

			totalItemPrice := item.Price

			itemTax := float64(totalItemPrice) / float64(subTotal) * float64(bill.Tax) / 100.0

			itemTotalWithTax := float64(totalItemPrice) + itemTax

			remaining := int(math.Round(itemTotalWithTax))
			numParticipants := len(participants)

			var itemParticipantResponses []dtos.GetBillByCreatorParticipantResponse
			if numParticipants > 0 {
				amountPerParticipant := remaining / numParticipants
				remainder := remaining % numParticipants

				for i, p := range participants {
					finalAmount := amountPerParticipant
					if i < remainder {
						finalAmount++
					}

					pResp := dtos.GetBillByCreatorParticipantResponse{
						ParticipantID: p.ParticipantID,
						AmountOwed:    float64(finalAmount) / 100.0,
						IsPaid:        p.IsPaid,
					}
					itemParticipantResponses = append(itemParticipantResponses, pResp)

					if existing, exists := participantMap[p.ParticipantID]; exists {
						existing.AmountOwed += float64(finalAmount) / 100.0
						participantMap[p.ParticipantID] = existing
					} else {
						participantMap[p.ParticipantID] = pResp
					}
				}
			}

			itemResponses = append(itemResponses, dtos.GetBillByCreatorItemResponse{
				ItemID:       item.ItemID,
				Name:         item.Name,
				Price:        float64(item.Price) / 100.0,
				Quantity:     item.Quantity,
				Participants: itemParticipantResponses,
			})
		}

		var participantResponses []dtos.GetBillByCreatorParticipantResponse
		for _, p := range participantMap {
			participantResponses = append(participantResponses, p)
		}

		response = append(response, dtos.GetBillByCreatorResponse{
			BillID:       bill.BillID,
			StoreName:    bill.StoreName,
			Tax:          float64(bill.Tax) / 100.0,
			CreatedAt:    bill.CreatedAt.Format(time.RFC3339),
			BillDate:     bill.BillDate,
			Items:        itemResponses,
			Participants: participantResponses,
		})
	}

	return response, nil
}

func AssignParticipantsToItem(req dtos.AssignParticipantsRequest) (*dtos.AssignedParticipantResponse, error) {
	var item models.Item
	if err := database.DB.First(&item, "item_id = ?", req.ItemID).Error; err != nil {
		return nil, errors.New("item not found")
	}

	if len(req.ParticipantID) == 0 {
		return nil, errors.New("no participant IDs provided")
	}

	var bill models.Bill
	if err := database.DB.First(&bill, "bill_id = ?", item.BillID).Error; err != nil {
		return nil, errors.New("bill not found for item")
	}

	var allItems []models.Item
	if err := database.DB.Where("bill_id = ?", bill.BillID).Find(&allItems).Error; err != nil {
		return nil, errors.New("failed to fetch items for tax calculation")
	}

	var subTotal int
	for _, i := range allItems {
		subTotal += i.Price
	}
	if subTotal == 0 {
		return nil, errors.New("subtotal is zero, cannot assign tax")
	}

	itemTax := (float64(item.Price) / float64(subTotal)) * float64(bill.Tax)
	itemTotalWithTax := float64(item.Price) + itemTax

	remaining := int(math.Round(itemTotalWithTax))
	numParticipants := len(req.ParticipantID)

	share := remaining / numParticipants
	remainder := remaining % numParticipants

	var assigned []dtos.AssignedParticipant

	for i, pid := range req.ParticipantID {
		amount := share
		if i < remainder {
			amount += 1
		}

		participant := models.Participant{
			ParticipantID: pid,
			ItemID:        item.ItemID,
			AmountOwed:    amount,
			IsPaid:        "",
		}

		if err := database.DB.Create(&participant).Error; err != nil {
			return nil, errors.New("failed to assign participant")
		}

		assigned = append(assigned, dtos.AssignedParticipant{
			ParticipantID: pid,
			ItemID:        item.ItemID,
			AmountOwed:    float64(amount) / 100.0,
			IsPaid:        false,
		})
	}

	resp := &dtos.AssignedParticipantResponse{
		ItemID:       item.ItemID,
		ItemName:     item.Name,
		Participants: assigned,
	}

	return resp, nil
}

func GetBillsByParticipantID(participantID string) ([]dtos.ParticipantBillResponse, error) {
	if participantID == "" {
		return nil, errors.New("participantId is required")
	}

	var participantRecords []models.Participant
	if err := database.DB.Where("participant_id = ?", participantID).Find(&participantRecords).Error; err != nil {
		return nil, errors.New("failed to fetch participant data: " + err.Error())
	}

	billMap := make(map[string]bool)
	for _, p := range participantRecords {
		var item models.Item
		if err := database.DB.Select("item_id", "bill_id").First(&item, "item_id = ?", p.ItemID).Error; err != nil {
			continue
		}
		billMap[item.BillID] = true
	}

	var responses []dtos.ParticipantBillResponse

	for billID := range billMap {
		var bill models.Bill
		if err := database.DB.First(&bill, "bill_id = ?", billID).Error; err != nil {
			continue
		}

		var items []models.Item
		if err := database.DB.Where("bill_id = ?", bill.BillID).Find(&items).Error; err != nil {
			continue
		}

		var itemDTOs []dtos.ParticipantItemResponse
		globalParticipantMap := make(map[string]int)

		subTotal := 0
		for _, item := range items {
			subTotal += item.Price
		}
		if subTotal == 0 {
			continue
		}

		for _, item := range items {
			var participants []models.Participant
			if err := database.DB.Where("item_id = ?", item.ItemID).Find(&participants).Error; err != nil {
				continue
			}

			totalItemPrice := item.Price

			itemTax := float64(totalItemPrice) / float64(subTotal) * float64(bill.Tax) / 100.0
			itemTotalWithTax := float64(totalItemPrice) + itemTax
			totalRounded := int(math.Round(itemTotalWithTax))

			numParticipants := len(participants)
			var participantResponses []dtos.ParticipantListResponse

			if numParticipants > 0 {
				amountPerParticipant := totalRounded / numParticipants
				remainder := totalRounded % numParticipants

				for i, p := range participants {
					finalAmount := amountPerParticipant
					if i < remainder {
						finalAmount++
					}

					participantResponses = append(participantResponses, dtos.ParticipantListResponse{
						ParticipantID: p.ParticipantID,
						AmountOwed:    float64(finalAmount) / 100.0,
						IsPaid:        p.IsPaid,
					})
					globalParticipantMap[p.ParticipantID] += finalAmount
				}
			}

			itemDTOs = append(itemDTOs, dtos.ParticipantItemResponse{
				ItemID:       item.ItemID,
				Name:         item.Name,
				Quantity:     item.Quantity,
				Price:        float64(item.Price) / 100.0,
				Participants: participantResponses,
			})
		}

		var globalParticipants []dtos.ParticipantListResponse
		for pid, totalAmount := range globalParticipantMap {
			var p models.Participant
			_ = database.DB.Where("participant_id = ?", pid).First(&p).Error

			globalParticipants = append(globalParticipants, dtos.ParticipantListResponse{
				ParticipantID: pid,
				AmountOwed:    float64(totalAmount) / 100.0,
				IsPaid:        p.IsPaid,
			})
		}

		responses = append(responses, dtos.ParticipantBillResponse{
			BillID:       bill.BillID,
			StoreName:    bill.StoreName,
			CreatorID:    bill.CreatorID,
			BillDate:     bill.BillDate,
			CreatedAt:    bill.CreatedAt.Format(time.RFC3339),
			Tax:          float64(bill.Tax) / 100.0,
			Items:        itemDTOs,
			Participants: globalParticipants,
		})
	}

	return responses, nil
}

func GetBillByBIllID(billID string) (dtos.ParticipantBillResponse, error) {
	if billID == "" {
		return dtos.ParticipantBillResponse{}, errors.New("billId is required")
	}

	var bill models.Bill
	if err := database.DB.Preload("Items").Where("bill_id = ?", billID).First(&bill).Error; err != nil {
		return dtos.ParticipantBillResponse{}, errors.New("bill not found")
	}

	var itemResponses []dtos.ParticipantItemResponse
	globalParticipantMap := make(map[string]int)

	var subTotal int
	for _, item := range bill.Items {
		subTotal += item.Price
	}
	if subTotal == 0 {
		return dtos.ParticipantBillResponse{}, errors.New("subtotal is 0, cannot calculate tax")
	}

	for _, item := range bill.Items {
		var participants []models.Participant
		if err := database.DB.Where("item_id = ?", item.ItemID).Find(&participants).Error; err != nil {
			return dtos.ParticipantBillResponse{}, errors.New("failed to get participants: " + err.Error())
		}

		totalItemPrice := item.Price
		itemTax := float64(totalItemPrice) / float64(subTotal) * float64(bill.Tax) / 100.0
		itemTotalWithTax := float64(totalItemPrice) + itemTax
		totalRounded := int(math.Round(itemTotalWithTax))

		numParticipants := len(participants)
		var participantResponses []dtos.ParticipantListResponse

		if numParticipants > 0 {
			amountPerParticipant := totalRounded / numParticipants
			remainder := totalRounded % numParticipants

			for i, p := range participants {
				finalAmount := amountPerParticipant
				if i < remainder {
					finalAmount++
				}

				participantResponses = append(participantResponses, dtos.ParticipantListResponse{
					ParticipantID: p.ParticipantID,
					AmountOwed:    float64(finalAmount) / 100.0,
					IsPaid:        p.IsPaid,
				})
				globalParticipantMap[p.ParticipantID] += finalAmount
			}
		}

		itemResponses = append(itemResponses, dtos.ParticipantItemResponse{
			ItemID:       item.ItemID,
			Name:         item.Name,
			Quantity:     item.Quantity,
			Price:        float64(item.Price) / 100.0,
			Participants: participantResponses,
		})
	}

	var globalParticipants []dtos.ParticipantListResponse
	for participantID, totalAmount := range globalParticipantMap {
		var p models.Participant
		_ = database.DB.Where("participant_id = ?", participantID).First(&p).Error

		globalParticipants = append(globalParticipants, dtos.ParticipantListResponse{
			ParticipantID: participantID,
			AmountOwed:    float64(totalAmount) / 100.0,
			IsPaid:        p.IsPaid,
		})
	}

	resp := dtos.ParticipantBillResponse{
		BillID:       bill.BillID,
		StoreName:    bill.StoreName,
		CreatorID:    bill.CreatorID,
		BillDate:     bill.BillDate,
		CreatedAt:    bill.CreatedAt.Format(time.RFC3339),
		Tax:          float64(bill.Tax) / 100.0,
		Items:        itemResponses,
		Participants: globalParticipants,
	}

	return resp, nil
}

func DeleteBillByIDService(billID string) (string, int, error) {
	if billID == "" {
		return "Bill ID is required", 400, errors.New("missing bill ID")
	}

	var bill models.Bill
	if err := database.DB.Where("bill_id = ?", billID).First(&bill).Error; err != nil {
		return "Bill not found", 404, err
	}

	var items []models.Item

	if err := database.DB.Where("bill_id = ?", billID).Find(&items).Error; err != nil {
		return "Item not found", 404, err
	}

	for _, item := range items {
		if err := database.DB.Where("item_id = ?", item.ItemID).Delete(&models.Participant{}).Error; err != nil {
			return "Failed to delete bill participants", 500, err
		}
	}

	if err := database.DB.Where("bill_id = ?", billID).Delete(&models.Item{}).Error; err != nil {
		return "Failed to delete bill items", 500, err
	}

	if err := database.DB.Delete(&bill).Error; err != nil {
		return "Failed to delete bill", 500, err
	}

	return "Delete bill success", 200, nil
}

func UpdateBillService(req dtos.UpdateBillRequest) (dtos.UpdateBillResponse, error) {
	var bill models.Bill
	if err := database.DB.Where("bill_id = ?", req.BillID).First(&bill).Error; err != nil {
		return dtos.UpdateBillResponse{}, err
	}

	var oldParticipants []models.Participant
	_ = database.DB.Table("participants").
		Joins("JOIN items ON participants.item_id = items.item_id").
		Where("items.bill_id = ?", bill.BillID).
		Find(&oldParticipants)

	existingParticipantMap := make(map[string]bool)
	for _, p := range oldParticipants {
		existingParticipantMap[p.ParticipantID] = true
	}

	var existingItems []models.Item
	if err := database.DB.Where("bill_id = ?", bill.BillID).Find(&existingItems).Error; err != nil {
		return dtos.UpdateBillResponse{}, err
	}
	for _, item := range existingItems {
		_ = database.DB.Where("item_id = ?", item.ItemID).Delete(&models.Participant{}).Error
	}
	_ = database.DB.Where("bill_id = ?", bill.BillID).Delete(&models.Item{}).Error

	bill.StoreName = req.StoreName
	bill.CreatorID = req.CreatorID
	bill.CreatedAt = time.Now()
	bill.BillDate = req.BillDate
	bill.Tax = utils.FormatUSDtoInt(req.Tax)
	if err := database.DB.Save(&bill).Error; err != nil {
		return dtos.UpdateBillResponse{}, err
	}

	var itemResponses []dtos.UpdateBillItemResponse
	subTotal := 0

	for _, item := range req.UpdateBillItemRequest {
		subTotal += utils.FormatUSDtoInt(item.Price)
	}

	notifiedParticipants := make(map[string]bool)
	paidNotificationSent := make(map[string]bool)

	for _, item := range req.UpdateBillItemRequest {
		var itemID string
		if err := database.DB.Where("item_id = ?", item.ItemID).First(&models.Item{}).Error; err == nil {
			itemID = item.ItemID
		} else {
			itemID = uuid.NewString()
		}

		priceInt := utils.FormatUSDtoInt(item.Price)
		totalItemPrice := priceInt

		itemTax := 0.0
		if subTotal > 0 {
			proportion := float64(totalItemPrice) / float64(subTotal)
			itemTax = proportion * float64(bill.Tax) / 100.9
		}
		itemTotalWithTax := float64(totalItemPrice) + itemTax
		itemTotalRounded := int(math.Round(itemTotalWithTax))

		newItem := models.Item{
			ItemID:   itemID,
			BillID:   bill.BillID,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    priceInt,
		}
		if err := database.DB.Save(&newItem).Error; err != nil {
			return dtos.UpdateBillResponse{}, err
		}

		var participantResponses []dtos.UpdateBillParticipantResponse
		participants := item.UpdateBillParticipantRequest
		numParticipants := len(participants)

		amountPerParticipant := 0
		remainder := 0
		if numParticipants > 0 {
			amountPerParticipant = itemTotalRounded / numParticipants
			remainder = itemTotalRounded % numParticipants
		}

		for i, p := range participants {
			amount := amountPerParticipant
			if i < remainder {
				amount++
			}

			newParticipant := models.Participant{
				ItemID:        newItem.ItemID,
				ParticipantID: p.ParticipantID,
				AmountOwed:    amount,
				IsPaid:        p.IsPaid,
			}
			if err := database.DB.Save(&newParticipant).Error; err != nil {
				return dtos.UpdateBillResponse{}, err
			}

			if p.ParticipantID != bill.CreatorID &&
				!existingParticipantMap[p.ParticipantID] &&
				p.IsPaid == "" &&
				!notifiedParticipants[p.ParticipantID] {

				handlers.SendInboxToUser(
					p.ParticipantID,
					fmt.Sprintf("You have been assigned in a bill: %s", bill.StoreName),
					"/participated-bills",
					"bill_assigned",
				)
				notifiedParticipants[p.ParticipantID] = true
			}

			if p.IsPaid != "" &&
				p.ParticipantID != bill.CreatorID &&
				!paidNotificationSent[p.ParticipantID] {

				handlers.SendInboxToUser(
					bill.CreatorID,
					fmt.Sprintf("Participant %s has paid their part in bill: %s", p.ParticipantID, bill.StoreName),
					"/created-bills",
					"participant_paid",
				)
				paidNotificationSent[p.ParticipantID] = true
			}

			participantResponses = append(participantResponses, dtos.UpdateBillParticipantResponse{
				ParticipantID: p.ParticipantID,
				AmountOwed:    float64(amount) / 100.0,
				IsPaid:        p.IsPaid,
			})
		}

		itemResponses = append(itemResponses, dtos.UpdateBillItemResponse{
			ItemID:                        newItem.ItemID,
			Name:                          newItem.Name,
			Quantity:                      newItem.Quantity,
			Price:                         float64(newItem.Price) / 100.0,
			UpdateBillParticipantResponse: participantResponses,
		})
	}

	return dtos.UpdateBillResponse{
		BillID:                 bill.BillID,
		StoreName:              bill.StoreName,
		CreatorID:              bill.CreatorID,
		CreatedAt:              bill.CreatedAt,
		BillDate:               bill.BillDate,
		Tax:                    float64(bill.Tax) / 100.0,
		UpdateBillItemResponse: itemResponses,
	}, nil
}
