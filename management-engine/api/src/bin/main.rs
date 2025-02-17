use api::auth::auth_routes;
use axum::{
    extract,
    http::StatusCode,
    response::{Html, IntoResponse, Response},
    routing::get,
    Router,
};

#[tokio::main]
async fn main() {
    // build our application with some routes
    let app = Router::new().nest("/auth", auth_routes());

    // run it
    let listener = tokio::net::TcpListener::bind("127.0.0.1:3000")
        .await
        .unwrap();
    axum::serve(listener, app).await.unwrap();
}
