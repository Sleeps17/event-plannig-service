package handlers

const (
	InternalErrorMsg  = "Something went wrong, please try again later"
	InvalidRequestMsg = "The request is not valid"
	InvalidIdParamMsg = "Param id is invalid"

	InvalidCredentialsMsg = "You have entered an incorrect username or password, please try again"
	UserAlreadyExistsMsg  = "A user with this username already exists"

	EmployeeNotFoundMsg   = "A employee with this id does not exist"
	EmployeeAlreadyExists = "A employee with such data already exists"

	RoomNotFoundMsg   = "A room with this id does not exist"
	RoomAlreadyExists = "A room with this id already exists"

	RoomIsOccupiedMsg       = "This room is occupied at this time"
	SomeEmployeesAreBusyMsg = "Some workers are busy at this time"
	EventNotFoundMsg        = "An event with this id does not exist"
)
