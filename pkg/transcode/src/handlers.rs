use actix_web::{web, HttpResponse, Result, post, get};
use serde_json::json;
use uuid::Uuid;
use chrono::Utc;

use crate::models::*;
use crate::services::*;
use crate::metadata::*;
use crate::filters::*;
use crate::ai::*;

#[get("/health")]
pub async fn health_check() -> Result<HttpResponse> {
    let response = HealthResponse {
        status: "ok".to_string(),
        message: "Transcode service is running".to_string(),
        timestamp: Utc::now(),
    };
    
    Ok(HttpResponse::Ok().json(response))
}

#[post("/video")]
pub async fn transcode_video(
    req: web::Json<VideoTranscodeRequest>,
) -> Result<HttpResponse> {
    let job_id = Uuid::new_v4();
    
    // TODO: Implement actual video transcode logic
    let response = TranscodeResponse {
        job_id,
        status: models::TranscodeStatus::Pending,
        message: "Video transcode job created".to_string(),
    };
    
    Ok(HttpResponse::Accepted().json(response))
}

#[post("/image")]
pub async fn transcode_image(
    req: web::Json<ImageTranscodeRequest>,
) -> Result<HttpResponse> {
    let job_id = Uuid::new_v4();
    
    // TODO: Implement actual image transcode logic
    let response = TranscodeResponse {
        job_id,
        status: models::TranscodeStatus::Pending,
        message: "Image transcode job created".to_string(),
    };
    
    Ok(HttpResponse::Accepted().json(response))
}

#[post("/audio")]
pub async fn transcode_audio(
    req: web::Json<AudioTranscodeRequest>,
) -> Result<HttpResponse> {
    let job_id = Uuid::new_v4();
    
    // TODO: Implement actual audio transcode logic
    let response = TranscodeResponse {
        job_id,
        status: models::TranscodeStatus::Pending,
        message: "Audio transcode job created".to_string(),
    };
    
    Ok(HttpResponse::Accepted().json(response))
}

#[get("/status/{job_id}")]
pub async fn get_transcode_status(
    path: web::Path<String>,
) -> Result<HttpResponse> {
    let job_id = path.into_inner();
    
    // TODO: Implement actual status checking logic
    let job = TranscodeJob {
        id: Uuid::parse_str(&job_id).unwrap_or_else(|_| Uuid::new_v4()),
        status: models::TranscodeStatus::Completed,
        input_path: "/tmp/input.jpg".to_string(),
        output_path: "/tmp/output.png".to_string(),
        format: "PNG".to_string(),
        progress: 100.0,
        error_message: None,
        created_at: Utc::now(),
        updated_at: Utc::now(),
    };
    
    Ok(HttpResponse::Ok().json(job))
}

// Metadata extraction endpoints
#[post("/extract")]
pub async fn extract_metadata(
    req: web::Json<MediaMetadataRequest>,
) -> Result<HttpResponse> {
    let extractor = MetadataExtractor::new();
    
    match extractor.extract_metadata(
        &req.file_path,
        req.extract_exif.unwrap_or(true),
        req.extract_ai_tags.unwrap_or(false)
    ).await {
        Ok(metadata) => Ok(HttpResponse::Ok().json(metadata)),
        Err(e) => Ok(HttpResponse::InternalServerError().json(json!({
            "error": "metadata_extraction_failed",
            "message": e.to_string()
        })))
    }
}

#[post("/analyze-video")]
pub async fn analyze_video(
    req: web::Json<VideoAnalysisRequest>,
) -> Result<HttpResponse> {
    let analyzer = VideoAnalyzer::new();
    
    match analyzer.analyze_video(
        &req.file_path,
        req.extract_frames.unwrap_or(false),
        req.frame_interval.unwrap_or(30),
        req.extract_audio.unwrap_or(false)
    ).await {
        Ok(analysis) => Ok(HttpResponse::Ok().json(analysis)),
        Err(e) => Ok(HttpResponse::InternalServerError().json(json!({
            "error": "video_analysis_failed",
            "message": e.to_string()
        })))
    }
}

// Image processing endpoints
#[post("/filter")]
pub async fn apply_filter(
    req: web::Json<ImageFilterRequest>,
) -> Result<HttpResponse> {
    let processor = ImageProcessor::new();
    
    match processor.apply_filter(&req).await {
        Ok(job_id) => Ok(HttpResponse::Accepted().json(TranscodeResponse {
            job_id,
            status: TranscodeStatus::Processing,
            message: "Filter application started".to_string(),
        })),
        Err(e) => Ok(HttpResponse::InternalServerError().json(json!({
            "error": "filter_application_failed",
            "message": e.to_string()
        })))
    }
}

#[post("/watermark")]
pub async fn add_watermark(
    req: web::Json<WatermarkRequest>,
) -> Result<HttpResponse> {
    let processor = ImageProcessor::new();
    
    match processor.add_watermark(&req).await {
        Ok(job_id) => Ok(HttpResponse::Accepted().json(TranscodeResponse {
            job_id,
            status: TranscodeStatus::Processing,
            message: "Watermark application started".to_string(),
        })),
        Err(e) => Ok(HttpResponse::InternalServerError().json(json!({
            "error": "watermark_failed",
            "message": e.to_string()
        })))
    }
}

#[post("/batch")]
pub async fn batch_process(
    req: web::Json<BatchProcessRequest>,
) -> Result<HttpResponse> {
    let processor = ImageProcessor::new();
    
    match processor.batch_process(&req).await {
        Ok(job_ids) => Ok(HttpResponse::Accepted().json(json!({
            "job_ids": job_ids,
            "message": format!("{} batch jobs started", job_ids.len())
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(json!({
            "error": "batch_processing_failed",
            "message": e.to_string()
        })))
    }
}

// AI analysis endpoints
#[post("/detect-objects")]
pub async fn detect_objects(
    req: web::Json<serde_json::Value>,
) -> Result<HttpResponse> {
    let image_path = req.get("image_path")
        .and_then(|v| v.as_str())
        .ok_or_else(|| actix_web::error::ErrorBadRequest("image_path is required"))?;
    
    let analyzer = AIAnalyzer::new();
    
    match analyzer.detect_objects(image_path).await {
        Ok(objects) => Ok(HttpResponse::Ok().json(json!({
            "objects": objects
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(json!({
            "error": "object_detection_failed",
            "message": e.to_string()
        })))
    }
}

#[post("/detect-faces")]
pub async fn detect_faces(
    req: web::Json<serde_json::Value>,
) -> Result<HttpResponse> {
    let image_path = req.get("image_path")
        .and_then(|v| v.as_str())
        .ok_or_else(|| actix_web::error::ErrorBadRequest("image_path is required"))?;
    
    let analyzer = AIAnalyzer::new();
    
    match analyzer.detect_faces(image_path).await {
        Ok(faces) => Ok(HttpResponse::Ok().json(json!({
            "faces": faces
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(json!({
            "error": "face_detection_failed",
            "message": e.to_string()
        })))
    }
}

#[post("/analyze-colors")]
pub async fn analyze_colors(
    req: web::Json<serde_json::Value>,
) -> Result<HttpResponse> {
    let image_path = req.get("image_path")
        .and_then(|v| v.as_str())
        .ok_or_else(|| actix_web::error::ErrorBadRequest("image_path is required"))?;
    
    let analyzer = AIAnalyzer::new();
    
    match analyzer.analyze_colors(image_path).await {
        Ok(colors) => Ok(HttpResponse::Ok().json(json!({
            "colors": colors
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(json!({
            "error": "color_analysis_failed",
            "message": e.to_string()
        })))
    }
} 