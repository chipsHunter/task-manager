use models::dto::user::UserGet;
use shared::response::ApiResponse;

pub struct UserController;

impl UserController {
    pub async fn login_user(conn: &mut sqlx::PgConnection, user: UserGet) -> ApiResponse {
        ApiResponse::Created
    }
}
