use clients::user::UserClient;
use models::dto::user::UserGet;
use shared::response::ApiResponse;

pub struct UserController;

impl UserController {
    pub async fn login_user(conn: &mut sqlx::PgConnection, user: UserGet) -> ApiResponse {
        let db_res = UserClient::get_user(user.clone().mail, conn).await;
        if let Ok(us) = db_res {
            if !us.nickname.is_empty() {
                return ApiResponse::BadRequest;
            }
        }
        match UserClient::create_user(user, conn).await {
            Ok(val) => ApiResponse::OK,
            Err(error) => ApiResponse::InternalServerError(error.to_string()),
        }
    }
}
