use models::{dto::user::UserGet, user::User};
use sqlx::types::Uuid;

pub struct UserClient;

impl UserClient {
    pub async fn get_user(
        mail: String,
        conn: &mut sqlx::PgConnection,
    ) -> Result<User, sqlx::Error> {
        sqlx::query_as!(
            User,
            "SELECT id, name, nickname, password, mail FROM users WHERE mail = $1",
            mail
        )
        .fetch_one(conn)
        .await
    }
    pub async fn create_user(
        user: UserGet,
        conn: &mut sqlx::PgConnection,
    ) -> Result<User, sqlx::Error> {
        sqlx::query_as!(
            User, // Тип, в который будем маппить результат
            "INSERT INTO users (id, name, nickname, password, mail) 
         VALUES ($1, $2, $3, $4, $5) 
         RETURNING id, name, nickname, password, mail", // Возвращаем поля для маппинга
            Uuid::new_v4(),
            user.name,
            user.nickname,
            user.password,
            user.mail
        )
        .fetch_one(conn)
        .await
    }
}
