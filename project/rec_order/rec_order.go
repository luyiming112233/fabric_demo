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
	Redeemed       = "Redeemed"
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
	OwnerID      string `json:"owner_id"`
	AcceptorID   string `json:"acceptor_id"`
	TotalAmount  int    `json:"total_amount"`
}

type Receivable struct {
	ReceivableNo        string `json:"receivable_no"`
	OrderNo             string `json:"order_no"`
	SignedTime          string `json:"signed_time"`
	ExpireTime          string `json:"expire_time"`
	OwnerID             string `json:"owner_id"`
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
			OwnerID:      "user2",
			AcceptorID:   "user1",
			TotalAmount:  1000000,
		},
		{
			OrderNo:      "order2",
			GoodsNo:      "goods2",
			ReceivableNo: "rec2",
			OwnerID:      "user2",
			AcceptorID:   "user1",
			TotalAmount:  1000000,
		},
		{
			OrderNo:      "order3",
			GoodsNo:      "goods3",
			ReceivableNo: "rec3",
			OwnerID:      "user2",
			AcceptorID:   "user1",
			TotalAmount:  1000000,
		},
		{
			OrderNo:      "order4",
			GoodsNo:      "goods4",
			ReceivableNo: "rec4",
			OwnerID:      "user2",
			AcceptorID:   "user1",
			TotalAmount:  500000,
		},
		{
			OrderNo:      "order5",
			GoodsNo:      "goods5",
			ReceivableNo: "rec5",
			OwnerID:      "user2",
			AcceptorID:   "user1",
			TotalAmount:  500000,
		},
		{
			OrderNo:      "order6",
			GoodsNo:      "goods6",
			ReceivableNo: "rec6",
			OwnerID:      "user2",
			AcceptorID:   "user1",
			TotalAmount:  500000,
		},
	}

	receivables := []Receivable{
		{
			ReceivableNo:        "rec2",
			OrderNo:             "order2",
			SignedTime:          "2020-05-27 15:04:05",
			ExpireTime:          "2021-05-27 15:04:05",
			OwnerID:             "user1",
			AcceptorID:          "user2",
			DiscountApplyAmount: 1000000,
			Status:              Accepted,
		},
		{
			ReceivableNo:        "rec3",
			OrderNo:             "order3",
			SignedTime:          "2020-05-27 15:04:05",
			ExpireTime:          "2021-05-27 15:04:05",
			OwnerID:             "user1",
			AcceptorID:          "user4",
			DiscountApplyAmount: 1000000,
		},
		{
			ReceivableNo:        "rec4",
			OrderNo:             "order4",
			SignedTime:          "2020-05-27 15:04:05",
			ExpireTime:          "2021-05-27 15:04:05",
			OwnerID:             "user1",
			AcceptorID:          "user2",
			DiscountApplyAmount: 500000,
			Status:              ToBeDiscounted,
		},
		{
			ReceivableNo:        "rec5",
			OrderNo:             "order5",
			SignedTime:          "2020-05-27 15:04:05",
			ExpireTime:          "2021-05-27 15:04:05",
			OwnerID:             "user1",
			AcceptorID:          "user4",
			DiscountApplyAmount: 500000,
			Status:              Discounted,
		},
		{
			ReceivableNo:        "rec6",
			OrderNo:             "order6",
			SignedTime:          "2020-05-27 15:04:05",
			ExpireTime:          "2020-05-28 15:04:05",
			OwnerID:             "user1",
			AcceptorID:          "user3",
			DiscountApplyAmount: 500000,
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
func (rc *ReceivableContract) CreateRecOrder(ctx contractapi.TransactionContextInterface, userID, acceptorID string, amount int) (*RecOrder, error) {
	// check userID is exist and the role is Company
	ok, err := rc.userCheck(ctx, userID, Company)
	if !ok {
		return nil, err
	}

	// check userID is exist and the role is Supplier
	ok, err = rc.userCheck(ctx, acceptorID, Supplier)
	if !ok {
		return nil, err
	}

	// First query all the receivable orders in order to determine the order number
	allOrders, err := rc.QueryAllRecOrders(ctx)
	if err != nil {
		return nil, err
	}
	orderCount := len(allOrders)
	orderCount++
	// OwnerID is Company, AcceptorID is Supplier
	order := &RecOrder{
		OrderNo:      "order" + strconv.Itoa(orderCount),
		GoodsNo:      "goods" + strconv.Itoa(orderCount),
		ReceivableNo: "rec" + strconv.Itoa(orderCount),
		OwnerID:      userID,
		AcceptorID:   acceptorID,
		TotalAmount:  amount,
	}

	dataBytes, _ := json.Marshal(order)
	err = ctx.GetStub().PutState("order"+strconv.Itoa(orderCount), dataBytes)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// SignReceivable using user2(First Supplier) to sign the receivable order
func (rc *ReceivableContract) SignReceivable(ctx contractapi.TransactionContextInterface, orderNo, supplierID string, discountApplyAmount int) (*Receivable, error) {
	// First query the receivable order by order number
	order, err := rc.QueryRecOrder(ctx, orderNo)
	if err != nil {
		return nil, err
	}

	// Then check the acceptor
	ok, err := rc.userCheck(ctx, order.OwnerID, Company)
	if !ok {
		return nil, err
	}

	// Check the signer
	ok, err = rc.userCheck(ctx, supplierID, Supplier)
	if !ok {
		return nil, err
	}

	// Check supplier is order.AcceptorID
	if supplierID != order.AcceptorID {
		return nil, fmt.Errorf("the supplierID %s is not equal with the AcceptorID %s in order", supplierID, order.AcceptorID)
	}

	// Check receivable whick the receivableNo is order.ReceivableNo is not exist
	if r, _ := rc.QueryReceivable(ctx, order.ReceivableNo); r != nil {
		return nil, fmt.Errorf("the Receivable which the receivableNo is %s exist", order.ReceivableNo)
	}

	// Check discountApplyAmount is less than or equal to order.TotalAmount
	if discountApplyAmount > order.TotalAmount {
		return nil, fmt.Errorf("discountApplyAmount %d is greater than totalAmount %d", discountApplyAmount, order.TotalAmount)
	}

	// Create the receivable proof
	rec := &Receivable{
		ReceivableNo:        order.ReceivableNo,
		OrderNo:             order.OrderNo,
		SignedTime:          time.Now().Format("2006-01-02 15:04:05"),
		ExpireTime:          time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05"),
		OwnerID:             supplierID,
		AcceptorID:          order.OwnerID,
		DiscountApplyAmount: discountApplyAmount,
		Status:              ToBeAccepted,
	}
	recBytes, _ := json.Marshal(rec)
	err = ctx.GetStub().PutState(rec.ReceivableNo, recBytes)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// AcceptRecOrder using user1(The Company) to accept the receivable
func (rc *ReceivableContract) AcceptReceivable(ctx contractapi.TransactionContextInterface, userID, receivableNo string) (*Receivable, error) {
	// Check userID is Company
	if ok, err := rc.userCheck(ctx, userID, Company); !ok {
		return nil, err
	}

	// Query the receivable by receivable number
	rec, err := rc.QueryReceivable(ctx, receivableNo)
	if err != nil {
		return nil, err
	}

	// Check the receivable's status is ToBeAccepted
	if ok, err := rc.receivableStatusCheck(rec, ToBeAccepted); !ok {
		return nil, err
	}

	rec.Status = Accepted
	recBytes, _ := json.Marshal(rec)
	err = ctx.GetStub().PutState(rec.ReceivableNo, recBytes)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// TransferReceivable transfers receivable from current acceptor(usually the First Supplier) to another acceptor(usually the Second Supplier)
func (rc *ReceivableContract) TransferReceivable(ctx contractapi.TransactionContextInterface, recNo, ownerID, supplierID string) (*Receivable, error) {
	// First query the receivable by receivable number
	rec, err := rc.QueryReceivable(ctx, recNo)
	if err != nil {
		return nil, err
	}

	// Check the receivable's status is Accepted
	if ok, err := rc.receivableStatusCheck(rec, Accepted); !ok {
		return nil, err
	}

	// Check the owner
	if ok, err := rc.userCheck(ctx, ownerID, Supplier); !ok {
		return nil, err
	}

	// Check the supplier
	if ok, err := rc.userCheck(ctx, supplierID, Supplier); !ok {
		return nil, err
	}

	// Check the rec.Owner
	if rec.OwnerID != ownerID {
		return nil, fmt.Errorf("user %s is not the owner of receivable %s, can't TransferReceivable", ownerID, recNo)
	}

	// Transfer the receivable, modify its acceptor
	rec.OwnerID = supplierID
	recBytes, _ := json.Marshal(rec)
	err = ctx.GetStub().PutState(rec.ReceivableNo, recBytes)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// ApplyDiscount is used by receivable acceptor(usually the Supplier) to apply discount from the financial organization
func (rc *ReceivableContract) ApplyDiscount(ctx contractapi.TransactionContextInterface, recNo, ownerID, financialID string) (*Receivable, error) {
	// First query the receivable by receivable number
	rec, err := rc.QueryReceivable(ctx, recNo)
	if err != nil {
		return nil, err
	}

	// Check the receivable status
	if ok, err := rc.receivableStatusCheck(rec, Accepted); !ok {
		return nil, err
	}

	// Check the owner
	if ok, err := rc.userCheck(ctx, ownerID, Supplier); !ok {
		return nil, err
	}

	// Check the financial
	if ok, err := rc.userCheck(ctx, financialID, Financial); !ok {
		return nil, err
	}

	// Check the rec.Owner
	if rec.OwnerID != ownerID {
		return nil, fmt.Errorf("user %s is not the owner of receivable %s, can't ApplyDiscount", ownerID, recNo)
	}

	// Apply for discount
	rec.Status = ToBeDiscounted
	recBytes, _ := json.Marshal(rec)
	err = ctx.GetStub().PutState(rec.ReceivableNo, recBytes)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func (rc *ReceivableContract) DiscountConfirm(ctx contractapi.TransactionContextInterface, recNo, financialID string) (*Receivable, error) {
	// First query the receivable by receivable number
	rec, err := rc.QueryReceivable(ctx, recNo)
	if err != nil {
		return nil, err
	}

	// Check the financial
	if ok, err := rc.userCheck(ctx, financialID, Financial); !ok {
		return nil, err
	}

	// Check the receivable status
	if ok, err := rc.receivableStatusCheck(rec, ToBeDiscounted); !ok {
		return nil, err
	}

	rec.Status = Discounted
	rec.OwnerID = financialID
	recBytes, _ := json.Marshal(rec)
	err = ctx.GetStub().PutState(rec.ReceivableNo, recBytes)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func (rc *ReceivableContract) Redeemed(ctx contractapi.TransactionContextInterface, recNo, companyID string) (*Receivable, error) {
	// First query the receivable by receivable number
	rec, err := rc.QueryReceivable(ctx, recNo)
	if err != nil {
		return nil, err
	}

	// Check the company
	if ok, err := rc.userCheck(ctx, companyID, Company); !ok {
		return nil, err
	}

	// Check the receivable status
	if ok, err := rc.receivableStatusCheck(rec, Discounted); !ok {
		return nil, err
	}

	// Check company is the acceptor
	if rec.AcceptorID != companyID {
		return nil, fmt.Errorf("user %s is not the acceptorID of receivable %s, can't Redeemed", companyID, rec.AcceptorID)
	}

	rec.Status = Redeemed
	rec.OwnerID = companyID
	recBytes, _ := json.Marshal(rec)
	err = ctx.GetStub().PutState(rec.ReceivableNo, recBytes)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func (rc *ReceivableContract) receivableStatusCheck(rec *Receivable, status string) (bool, error) {
	// Check the receivable's status
	if rec.Status != status {
		return false, fmt.Errorf("the receivable's status is %s while %s is needed", rec.Status, status)
	}

	return true, nil
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
