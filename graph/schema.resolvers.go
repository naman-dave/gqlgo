package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/naman-dave/gqlgo/graph/generated"
	"github.com/naman-dave/gqlgo/graph/model"
	"github.com/naman-dave/gqlgo/internal/jwt"
	"github.com/naman-dave/gqlgo/internal/middleware"
	dbmodel "github.com/naman-dave/gqlgo/internal/model"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	user := dbmodel.User{}
	user.Firstname = input.Firstname
	user.Lastname = input.Lastname
	user.Mobilenumber = input.Mobilenumber
	user.Passkey = input.Password

	err := user.Create()
	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(user.Mobilenumber)
	if err != nil {
		return "", err
	}

	user.Token = token
	err = user.UpdateToken()
	if err != nil {
		return "", err
	}
	return token, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user dbmodel.User
	user.Mobilenumber = input.Mobilenumber
	user.Passkey = input.Password

	correct, err := user.Authenticate()
	if err != nil {
		return "", fmt.Errorf("user does not exist")
	}

	if !correct {
		return "", fmt.Errorf("token is invalid, please login again")
	}

	token, err := jwt.GenerateToken(user.Mobilenumber)
	if err != nil {
		return "", err
	}

	user.Token = token
	err = user.UpdateToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

// AddCar is the resolver for the addCar field.
func (r *mutationResolver) AddCar(ctx context.Context, input model.NewCar) (string, error) {
	_, err := middleware.CheckIfLoggedIn(ctx)
	if err != nil {
		return "", err
	}

	_, err = time.Parse("2006-01-02", input.Dateofmanufacture)
	if err != nil {
		return "", fmt.Errorf(dateParseErr, "dateofmanufacture")
	}

	car := dbmodel.Car{
		CarIdentifier:     input.Caridentifier,
		Name:              input.Modal,
		DateOfManufacture: input.Dateofmanufacture,
		TotalCar:          input.Totalcar,
		TotalInUse:        input.Totalinuse,
	}

	id, err := car.Insert()
	if err != nil {
		return "", fmt.Errorf("can not add car, please try again!")
	}

	return id, nil
}

// BookCar is the resolver for the bookCar field.
func (r *mutationResolver) BookCar(ctx context.Context, input model.ProcessCar) (string, error) {
	loginuserID, err := middleware.CheckIfLoggedIn(ctx)
	if err != nil {
		return "", err
	}

	_, err = time.Parse("2006-01-02", input.Bookedtill)
	if err != nil {
		return "", fmt.Errorf(dateParseErr, "booktill")
	}

	booking := dbmodel.CarUsage{}

	booking.CarUniqueID = input.Caridentifier
	booking.BookedTill = input.Bookedtill
	booking.UserID = strconv.Itoa(loginuserID)

	bookingNO, err := booking.BookCar()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(bookingNO)), nil
}

// ReturnCar is the resolver for the returnCar field.
func (r *mutationResolver) ReturnCar(ctx context.Context, input int) (string, error) {
	_, err := middleware.CheckIfLoggedIn(ctx)
	if err != nil {
		return "", err
	}

	booking, err := dbmodel.GetCarUsage(int64(input))
	if err != nil {
		return "", fmt.Errorf("can not find booking for %d", input)
	}

	if booking.RetunedDate != "" {
		return "", fmt.Errorf("car is already returned on: %s", booking.RetunedDate)
	}

	booking.RetunedDate = time.Now().Format("2006-01-02")

	err = booking.ReturnCar(int64(input))
	if err != nil {
		return "", err
	}

	return "Returned", nil
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context) (string, error) {
	userID, err := middleware.CheckIfLoggedIn(ctx)
	if err != nil {
		return "", err
	}

	err = dbmodel.RemoveToken(userID)
	if err != nil {
		return "", err
	}

	return "Logged out", nil
}

// Cars is the resolver for the cars field.
func (r *queryResolver) Cars(ctx context.Context) ([]*model.Car, error) {
	_, err := middleware.CheckIfLoggedIn(ctx)
	if err != nil {
		return nil, err
	}

	cars := []*model.Car{}

	err = dbmodel.GetCars(&cars)
	if err != nil {
		return nil, err
	}

	return cars, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
const (
	dateParseErr = "%s is not in valid format. valid format (YYYY-MM-DD)"
)
