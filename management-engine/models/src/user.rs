use uuid::Uuid;

pub struct User {
    pub id: Uuid,
    pub name: Option<String>,
    pub nickname: String,
    pub mail: String,
    pub password: String,
}
