use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use serde::Serialize;

// here we show a type that implements Serialize + Send
#[derive(Serialize)]
pub struct Message {
    message: String,
}

pub enum ApiResponse {
    OK,
    Created,
    BadRequest,
    JsonData(Vec<Message>),
}

impl IntoResponse for ApiResponse {
    fn into_response(self) -> Response {
        match self {
            Self::OK => (StatusCode::OK).into_response(),
            Self::Created => (StatusCode::CREATED).into_response(),
            Self::BadRequest => (StatusCode::BAD_REQUEST).into_response(),
            Self::JsonData(data) => (StatusCode::OK, Json(data)).into_response(),
        }
    }
}
