use axum::{
    extract::State,
    http::StatusCode,
    response::IntoResponse,
    routing::{get, post},
    Json, Router,
};
use controllers::user::UserController;
use models::{dto::user::UserGet, user::User};
use serde::{Deserialize, Serialize};
use shared::response::{ApiResponse, Message};
use sqlx::PgPool;

pub async fn login(State(pool): State<PgPool>, Json(payload): Json<UserGet>) -> impl IntoResponse {
    let mut conn = pool.acquire().await.unwrap();
    UserController::login_user(&mut conn, payload).await
}
