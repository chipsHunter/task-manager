use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Default, JsonSchema)]
pub struct User {
    pub name: String,
    pub nickname: String,
    pub mail: String,
    pub password: String,
}
