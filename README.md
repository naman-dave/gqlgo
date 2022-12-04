Here, I was not able to get how we are going to make an microservice with two different services with use of the graphql. means if we add two services with graphql it doesn't make sense. that's why created the monolith.

run `docker-compose up`

run `migrate -database postgres://postgres:postgres@localhost:5432?sslmode=disable -path internal/db/migrations up`

```
User {
    id
    firstname
    lastname
    mobilenumber
    passkey
}

Method(User) 
 - login
 - logout
 - signin

Car {
    id
    model
    caridentifier
    dateofmenufactore
    totalcars
    totalcarinuse
}

Method(Car)
  - getallcar
  - addcar

CarUsgae {
   id (billingid)
   caridentifier FK (Car.caridentifier)
   userid FK (User.id)
   bookedtill
   returndate
}

Method(CarUsage)
  - bookcar
  - returncar
```
