use models::user::User;

async fn get_user(mail: String, conn: &mut sqlx::PgConnection) -> Result<User, sqlx::Error> {
    sqlx::query_as!(
        User,
        "SELECT id, name, nickname, password, mail FROM users WHERE mail = $1",
        mail
    )
    .fetch_one(conn)
    .await
}
