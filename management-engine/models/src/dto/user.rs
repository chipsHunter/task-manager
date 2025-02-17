use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Default, JsonSchema, Clone, Debug)]
pub struct UserGet {
    pub name: String,
    pub nickname: String,
    pub mail: String,
    pub password: String,
}
