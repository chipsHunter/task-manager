use axum::{
    routing::{get, post},
    Json, Router,
};
use models::user::User;
use serde::{Deserialize, Serialize};
use shared::response::ApiResponse;

pub fn auth_routes() -> Router {
    Router::new().route("/login", post(login)) // Добавляем маршрут логина
}

async fn login(Json(payload): Json<User>) -> ApiResponse {
    // Здесь могла быть проверка логина, но пока просто возвращаем фиктивный токен
    ApiResponse::Created
}
