use std::{env, time::Duration};

use api::auth::login;
use axum::{
    extract::{FromRef, FromRequestParts, State},
    http::{
        header::{AUTHORIZATION, CONTENT_TYPE},
        request::Parts,
        Method, StatusCode,
    },
    routing::{get, post},
    Router,
};
use dotenv::dotenv;
use sqlx::postgres::{PgPool, PgPoolOptions};
use tokio::net::TcpListener;
use tower_http::cors::{AllowOrigin, Any, CorsLayer};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

// we can extract the connection pool with `State`
async fn using_connection_pool_extractor(
    State(pool): State<PgPool>,
) -> Result<String, (StatusCode, String)> {
    sqlx::query_scalar("select 'hello world from pg'")
        .fetch_one(&pool)
        .await
        .map_err(internal_error)
}

fn internal_error<E>(err: E) -> (StatusCode, String)
where
    E: std::error::Error,
{
    (StatusCode::INTERNAL_SERVER_ERROR, err.to_string())
}

#[tokio::main]
async fn main() {
    dotenv().ok();

    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| format!("{}=debug", env!("CARGO_CRATE_NAME")).into()),
        )
        .with(tracing_subscriber::fmt::layer())
        .init();

    let cors = CorsLayer::new()
        .allow_origin(Any) // Указываем конкретный источник
        .allow_methods([
            Method::GET,
            Method::POST,
            Method::PUT,
            Method::DELETE,
            Method::OPTIONS,
        ]) // Разрешенные HTTP-методы
        .allow_headers([AUTHORIZATION, CONTENT_TYPE]) // Разрешенные заголовки
        //.allow_credentials(true) // Разрешаем отправку credentials (cookies, авторизация)
        .max_age(Duration::from_secs(3600)); // Кэшируем CORS-настройки

    let pool = PgPoolOptions::new()
        .max_connections(5)
        .acquire_timeout(Duration::from_secs(3))
        .connect(&env::var("DATABASE_URL").unwrap())
        .await
        .expect("can't connect to database");

    // build our application with some routes
    let app = Router::new()
        .route("/login", post(login).options(|| async { StatusCode::OK }))
        .route(
            "/",
            get(using_connection_pool_extractor).post(using_connection_pool_extractor),
        )
        .with_state(pool)
        .layer(cors);

    // run it
    let listener = tokio::net::TcpListener::bind("0.0.0.0:4000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
