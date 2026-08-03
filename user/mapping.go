package user

func toUser(id ID, request *Request) *User {
	return &User{
		ID:       id,
		Name:     request.Name,
		Password: request.Password,
	}
}

func toResponse(user *User) *Response {
	return &Response{
		ID:   user.ID,
		Name: user.Name,
	}
}

func toResponses(users []User) []Response {
	responses := make([]Response, len(users))

	for i, user := range users {
		responses[i] = *toResponse(&user)
	}

	return responses
}

func toIDResponse(id ID) *IDResponse {
	return &IDResponse{
		ID: id,
	}
}
