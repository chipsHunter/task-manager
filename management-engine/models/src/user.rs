use uuid::Uuid;

#[derive(Debug)]
pub struct User {
    pub id: Uuid,
    pub name: Option<String>,
    pub nickname: String,
    pub mail: String,
    pub password: String,
}
