package response

// 403 {"error":"API key is missing."}
// 400 {"error":"Credit card data invalid."}
// 400 {"error":"Please provide *** fields."}
// 401 {"error":"Invalid API key."}
// 500 {"error":"Something went wrong. Please try again later."}

type ErrorMessage string
type ErrorCode int

const (
	MISSING_API_MSG         ErrorMessage = "API key is missing."
	INVALID_CREDIT_CARD_MSG ErrorMessage = "Credit card data invalid."
	MISSING_FIELD_MSG       ErrorMessage = "Please provide %v fields."
	INVALID_API_KEY_MSG     ErrorMessage = "Invalid API key."
	INTERNAL_ERROR_MSG      ErrorMessage = "Something went wrong. Please try again later."
	BAD_REQUEST_MSG         ErrorMessage = "Invalid Requested Data. Please Contact Our Admin."

	MISSING_API_CODE         ErrorCode = 403
	INVALID_CREDIT_CARD_CODE ErrorCode = 400
	MISSING_FIELD_CODE       ErrorCode = 400
	BAD_REQUEST_CODE         ErrorCode = 400
	INVALID_API_KEY_CODE     ErrorCode = 401
	INTERNAL_ERROR_CODE      ErrorCode = 500
)

type ErrorResponse struct {
	Code  ErrorCode    `json:"-"`
	Error ErrorMessage `json:"error"`
}
