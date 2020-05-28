package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// User Roles
const (
	Company   = "Company"
	Supplier  = "Supplier"
	Financial = "Financial"
)

// Account Status
const (
	Invalid = "Invalid"
	Valid   = "Valid"
	Frozen  = "Frozen"
)

// Receivable Order Status
const (
	ToBeAccepted   = "To Be Accepted"
	Accepted       = "Accepted"
	ToBeDiscounted = "To Be Discounted"
	Discounted     = "Discounted"
	Paid           = "Paid"
	Expired        = "Expired"
)

type Account struct {
	ID           string `json:"id"`
	Enterprise   string `json:"enterprise"`
	Role         string `json:"role"`
	Status       string `json:"status"`
	CertNo       string `json:"cert_no"`
	AcctSvcrName string `json:"acct_svcr_name"`
}

type RecOrder struct {
	OrderNo      string `json:"order_no"`
	GoodsNo      string `json:"goods_no"`
	ReceivableNo string `json:"receivable_no"`
	SingerID     string `json:"singer_id"`
	AcceptorID   string `json:"acceptor_id"`
	TotalAmount  int    `json:"total_amount"`
	Status       string `json:"status"`
}

type Receivable struct {
	ReceivableNo        string `json:"receivable_no"`
	OrderNo             string `json:"order_no"`
	SignedTime          string `json:"signed_time"`
	ExpireTime          string `json:"expire_time"`
	SingerID            string `json:"singer_id"`
	AcceptorID          string `json:"acceptor_id"`
	DiscountApplyAmount int    `json:"discount_apply_amount"`
	Status              string `json:"status"`
}

type QueryResult struct {
	Key    string `json:"key"`
	Record interface{}
}

type ReceivableContract struct {
	contractapi.Contract
}

func (rc *ReceivableContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	users := []Account{
		{
			ID:           "user1",
			Enterprise:   "Qulian Technology",
			Role:         Company,
			Status:       Valid,
			CertNo:       "cert number for user1",
			AcctSvcrName: "account service bank name for user1",
		},
		{
			ID:           "user2",
			Enterprise:   "First Supplier",
			Role:         Supplier,
			Status:       Valid,
			CertNo:       "cert number for user2",
			AcctSvcrName: "account service bank name for user2",
		},
		{
			ID:           "user3",
			Enterprise:   "Second Supplier",
			Role:         Supplier,
			Status:       Valid,
			CertNo:       "cert number for user3",
			AcctSvcrName: "account service bank name for user3",
		},
		{
			ID:           "user4",
			Enterprise:   "XX Bank",
			Role:         Financial,
			Status:       Valid,
			CertNo:       "cert number for user4",
			AcctSvcrName: "account service bank name for user4",
		},
	}

	recOrders := []RecOrder{
		{
			OrderNo:      "order1",
			GoodsNo:      "goods1",
			ReceivableNo: "rec1",
			SingerID:     "user2",
			AcceptorID:   "user1",
			TotalAmount:  1000000,
			Status:       ToBeAccepted,
		},
		{
			OrderNo:      "order2",
			GoodsNo:      "goods2",
			ReceivableNo: "rec2",
			SingerID:     "user2",
			AcceptorID:   "user1",
			TotalAmount:  1000000,
			Status:       Accepted,
		},
		{
			OrderNo:      "order3",
			GoodsNo:      "goods3",
			ReceivableNo: "rec3",
			SingerID:     "user2",
			AcceptorID:   "user1",
			TotalAmount:  1000000,
			Status:       Paid,
		},
		{
			OrderNo:      "order4",
			GoodsNo:      "goods4",
			ReceivableNo: "rec4",
			SingerID:     "user2",
			AcceptorID:   "user1",
			TotalAmount:  500000,
			Status:       ToBeDiscounted,
		},
		{
			OrderNo:      "order5",
			GoodsNo:      "goods5",
			ReceivableNo: "rec5",
			SingerID:     "user2",
			AcceptorID:   "user1",
			TotalAmount:  500000,
			Status:       Discounted,
		},
		{
			OrderNo:      "order6",
			GoodsNo:      "goods6",
			ReceivableNo: "rec6",
			SingerID:     "user2",
			AcceptorID:   "user1",
			TotalAmount:  500000,
			Status:       Expired,
		},
	}

	receivables := []Receivable{
		{
			ReceivableNo:        "rec2",
			OrderNo:             "order2",
			SignedTime:          "2020-05-27 15:04:05",
			ExpireTime:          "2021-05-27 15:04:05",
			SingerID:            "user1",
			AcceptorID:          "user2",
			DiscountApplyAmount: 1000000,
			Status:              Accepted,
		},
		{
			ReceivableNo:        "rec3",
			OrderNo:             "order3",
			SignedTime:          "2020-05-27 15:04:05",
			ExpireTime:          "2021-05-27 15:04:05",
			SingerID:            "user1",
			AcceptorID:          "user4",
			DiscountApplyAmount: 1000000,
			Status:              Paid,
		},
		{
			ReceivableNo:        "rec4",
			OrderNo:             "order4",
			SignedTime:          "2020-05-27 15:04:05",
			ExpireTime:          "2021-05-27 15:04:05",
			SingerID:            "user1",
			AcceptorID:          "user2",
			DiscountApplyAmount: 500000,
			Status:              ToBeDiscounted,
		},
		{
			ReceivableNo:        "rec5",
			OrderNo:             "order5",
			SignedTime:          "2020-05-27 15:04:05",
			ExpireTime:          "2021-05-27 15:04:05",
			SingerID:            "user1",
			AcceptorID:          "user4",
			DiscountApplyAmount: 500000,
			Status:              Discounted,
		},
		{
			ReceivableNo:        "rec6",
			OrderNo:             "order6",
			SignedTime:          "2020-05-27 15:04:05",
			ExpireTime:          "2020-05-28 15:04:05",
			SingerID:            "user1",
			AcceptorID:          "user3",
			DiscountApplyAmount: 500000,
			Status:              Expired,
		},
	}

	for _, user := range users {
		id := user.ID
		userAsBytes, err := json.Marshal(user)
		if err != nil {
			return fmt.Errorf("unable to marshal user account object, %s", err.Error())
		}
		err = ctx.GetStub().PutState(id, userAsBytes)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %s", err.Error())
		}
	}

	for _, recOrder := range recOrders {
		orderNo := recOrder.OrderNo
		orderAsBytes, err := json.Marshal(recOrder)
		if err != nil {
			return fmt.Errorf("unable to marshal receivable order object, %s", err.Error())
		}
		err = ctx.GetStub().PutState(orderNo, orderAsBytes)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %s", err.Error())
		}
	}

	for _, rec := range receivables {
		recNo := rec.ReceivableNo
		recAsBytes, err := json.Marshal(rec)
		if err != nil {
			return fmt.Errorf("unable to marshal receivable object, %s", err.Error())
		}
		err = ctx.GetStub().PutState(recNo, recAsBytes)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

func (rc *ReceivableContract) QueryUser(ctx contractapi.TransactionContextInterface, userID string) (*Account, error) {
	dataBytes, err := ctx.GetStub().GetState(userID)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	if dataBytes == nil {
		return nil, fmt.Errorf("%s does not exist", userID)
	}

	user := new(Account)
	err = json.Unmarshal(dataBytes, user)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal user data bytes, %s", err.Error())
	}
	return user, nil
}

func (rc *ReceivableContract) QueryAllUsers(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "user0"
	endKey := "user99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		user := new(Account)
		_ = json.Unmarshal(queryResponse.Value, user)

		queryResult := QueryResult{Key: queryResponse.Key, Record: user}
		results = append(results, queryResult)
	}

	return results, nil
}

func (rc *ReceivableContract) QueryRecOrder(ctx contractapi.TransactionContextInterface, orderNo string) (*RecOrder, error) {
	dataBytes, err := ctx.GetStub().GetState(orderNo)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	if dataBytes == nil {
		return nil, fmt.Errorf("%s does not exist", orderNo)
	}

	order := new(RecOrder)
	err = json.Unmarshal(dataBytes, order)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal user data bytes, %s", err.Error())
	}
	return order, nil
}

func (rc *ReceivableContract) QueryAllRecOrders(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "order0"
	endKey := "order99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		order := new(RecOrder)
		_ = json.Unmarshal(queryResponse.Value, order)

		queryResult := QueryResult{Key: queryResponse.Key, Record: order}
		results = append(results, queryResult)
	}

	return results, nil
}

func (rc *ReceivableContract) QueryReceivable(ctx contractapi.TransactionContextInterface, recNo string) (*Receivable, error) {
	dataBytes, err := ctx.GetStub().GetState(recNo)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	if dataBytes == nil {
		return nil, fmt.Errorf("%s does not exist", recNo)
	}

	rec := new(Receivable)
	err = json.Unmarshal(dataBytes, rec)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal user data bytes, %s", err.Error())
	}
	return rec, nil
}

func (rc *ReceivableContract) QueryAllReceivables(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "rec0"
	endKey := "rec99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		rec := new(Receivable)
		_ = json.Unmarshal(queryResponse.Value, rec)

		queryResult := QueryResult{Key: queryResponse.Key, Record: rec}
		results = append(results, queryResult)
	}

	return results, nil
}

func (rc *ReceivableContract) userCheck(ctx contractapi.TransactionContextInterface, userID, role string) (bool, error) {
	// Check whether the user exists
	user, err := rc.QueryUser(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("no such user account: %s", userID)
	}

	// Check the user's role
	if user.Role != role {
		return false, fmt.Errorf("the user is a %s while a %s is needed", user.Role, role)
	}

	return true, nil
}

// CreateRecOrder using user1(The Company) to create a receivable order
func (rc *ReceivableContract) CreateRecOrder(ctx contractapi.TransactionContextInterface, amount int) (*RecOrder, error) {
	// First query all the receivable orders in order to determine the order number
	allOrders, err := rc.QueryAllRecOrders(ctx)
	if err != nil {
		return nil, err
	}
	orderCount := len(allOrders)
	orderCount++
	order := &RecOrder{
		OrderNo:      "order" + strconv.Itoa(orderCount),
		GoodsNo:      "goods" + strconv.Itoa(orderCount),
		ReceivableNo: "rec" + strconv.Itoa(orderCount),
		SingerID:     "",
		AcceptorID:   "user1",
		TotalAmount:  amount,
		Status:       ToBeAccepted,
	}

	dataBytes, _ := json.Marshal(order)
	err = ctx.GetStub().PutState("order"+strconv.Itoa(orderCount), dataBytes)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// SignRecOrder using user2(First Supplier) to sign the receivable order
func (rc *ReceivableContract) SignRecOrder(ctx contractapi.TransactionContextInterface, orderNo, supplierID string) (*RecOrder, error) {
	// First query the receivable order by order number
	order, err := rc.QueryRecOrder(ctx, orderNo)
	if err != nil {
		return nil, err
	}

	// Then check the acceptor
	ok, err := rc.userCheck(ctx, order.AcceptorID, Company)
	if !ok {
		return nil, err
	}

	// Check the signer
	ok, err = rc.userCheck(ctx, supplierID, Supplier)
	if !ok {
		return nil, err
	}

	// Sign the order
	order.SingerID = supplierID
	dataBytes, _ := json.Marshal(order)
	err = ctx.GetStub().PutState(order.OrderNo, dataBytes)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// AcceptRecOrder using user1(The Company) to accept the receivable
func (rc *ReceivableContract) AcceptRecOrder(ctx contractapi.TransactionContextInterface, orderNo string, discount int) (*Receivable, error) {
	// First query the receivable order by order number
	order, err := rc.QueryRecOrder(ctx, orderNo)
	if err != nil {
		return nil, err
	}

	// Then accept the receivable
	order.Status = Accepted
	orderBytes, _ := json.Marshal(order)
	err = ctx.GetStub().PutState(order.OrderNo, orderBytes)
	if err != nil {
		return nil, err
	}

	// If the discount apply amount is larger than total amount, drop it
	if discount > order.TotalAmount {
		return nil, fmt.Errorf("discount apply amount is larger than receivable order's total amount")
	}

	// Create the receivable proof
	rec := &Receivable{
		ReceivableNo:        order.ReceivableNo,
		OrderNo:             order.OrderNo,
		SignedTime:          time.Now().Format("2006-01-02 15:04:05"),
		ExpireTime:          time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05"),
		SingerID:            "user1",
		AcceptorID:          order.SingerID,
		DiscountApplyAmount: discount,
		Status:              Accepted,
	}
	recBytes, _ := json.Marshal(rec)
	err = ctx.GetStub().PutState(rec.ReceivableNo, recBytes)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// TransferReceivable transfers receivable from current acceptor(usually the First Supplier) to another acceptor(usually the Second Supplier)
func (rc *ReceivableContract) TransferReceivable(ctx contractapi.TransactionContextInterface, recNo, supplierID string) (*Receivable, error) {
	// First query the receivable by receivable number
	rec, err := rc.QueryReceivable(ctx, recNo)
	if err != nil {
		return nil, err
	}
	
	// Check the receivable status
	if rec.Status != Accepted {
		return nil, fmt.Errorf("the receivable is not accepted")
	}

	// Check the supplier
	ok, err := rc.userCheck(ctx, supplierID, Supplier)
	if !ok {
		return nil, err
	}

	// Transfer the receivable, modify its acceptor
	rec.AcceptorID = supplierID
	recBytes, _ := json.Marshal(rec)
	err = ctx.GetStub().PutState(rec.ReceivableNo, recBytes)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// ApplyDiscount is used by receivable acceptor(usually the Supplier) to apply discount from the financial organization
func (rc *ReceivableContract) ApplyDiscount(ctx contractapi.TransactionContextInterface, recNo, financialID string) (*Receivable, error) {
	// First query the receivable by receivable number
	rec, err := rc.QueryReceivable(ctx, recNo)
	if err != nil {
		return nil, err
	}

	// Check the receivable status
	if rec.Status != Accepted {
		return nil, fmt.Errorf("the receivable is not accepted")
	}

	// Check the financial
	ok, err := rc.userCheck(ctx, financialID, Financial)
	if !ok {
		return nil, err
	}

	// Apply for discount
	rec.Status = ToBeDiscounted
	recBytes, _ := json.Marshal(rec)
	err = ctx.GetStub().PutState(rec.ReceivableNo, recBytes)
	if err != nil {
		return nil, err
	}

	// Make sure the order is updated consistently
	order, err := rc.QueryRecOrder(ctx, rec.OrderNo)
	order.Status = ToBeDiscounted
	orderBytes, _ := json.Marshal(order)
	err = ctx.GetStub().PutState(order.OrderNo, orderBytes)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func (rc *ReceivableContract) ConfirmDiscountApplication(ctx contractapi.TransactionContextInterface, orderNo, financialID string) (*Receivable, error) {
	// First query the receivable order by order number
	order, err := rc.QueryRecOrder(ctx, orderNo)
	if err != nil {
		return nil, err
	}

	// Check the receivable status
	if order.Status != ToBeDiscounted {
		return nil, fmt.Errorf("the receivable has not been applied for discount yet")
	}

	// Check the receivable's signer(also the receivable order acceptor)
	ok, err := rc.userCheck(ctx, order.AcceptorID, Company)
	if !ok {
		return nil, err
	}

	// Check the receivable order's signer
	ok, err = rc.userCheck(ctx, order.SingerID, Supplier)
	if !ok {
		return nil, err
	}

	// Query for the receivable
	rec, err := rc.QueryReceivable(ctx, order.ReceivableNo)
	if err != nil {
		return nil, err
	}

	// Check the receivable's acceptor(the receivable is possibly transferred, so the acceptor can be another supplier)
	ok, err = rc.userCheck(ctx, rec.AcceptorID, Supplier)
	if err != nil {
		return nil, err
	}

	// Check the financial that is going to be the receivable's acceptor
	ok, err = rc.userCheck(ctx, financialID, Financial)
	if err != nil {
		return nil, err
	}

	// Now do the discount
	rec.Status = Discounted
	rec.AcceptorID = financialID
	recBytes, _ := json.Marshal(rec)
	err = ctx.GetStub().PutState(rec.ReceivableNo, recBytes)
	if err != nil {
		return nil, err
	}

	order.Status = Discounted
	orderBytes, _ := json.Marshal(order)
	err = ctx.GetStub().PutState(order.OrderNo, orderBytes)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(ReceivableContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
