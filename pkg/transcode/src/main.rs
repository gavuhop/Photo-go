use actix_web::{web, App, HttpServer, middleware::Logger};
use log::info;
use std::env;

mod handlers;
mod models;
mod services;
mod error;
mod metadata;
mod filters;
mod ai;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));

    let port = env::var("PORT").unwrap_or_else(|_| "8081".to_string());
    let addr = format!("127.0.0.1:{}", port);

    info!("Starting transcode service on {}", addr);

    HttpServer::new(|| {
        App::new()
            .wrap(Logger::default())
            .service(
                web::scope("/api/v1")
                    .service(handlers::health_check)
                    .service(
                        web::scope("/transcode")
                            .service(handlers::transcode_video)
                            .service(handlers::transcode_image)
                            .service(handlers::transcode_audio)
                            .service(handlers::get_transcode_status)
                    )
                    .service(
                        web::scope("/metadata")
                            .service(handlers::extract_metadata)
                            .service(handlers::analyze_video)
                    )
                    .service(
                        web::scope("/process")
                            .service(handlers::apply_filter)
                            .service(handlers::add_watermark)
                            .service(handlers::batch_process)
                    )
                    .service(
                        web::scope("/ai")
                            .service(handlers::detect_objects)
                            .service(handlers::detect_faces)
                            .service(handlers::analyze_colors)
                    )
            )
    })
    .bind(addr)?
    .run()
    .await
} 