use axum::{
    Router,
    extract::State,
    http::StatusCode,
    routing::{delete_service, get, post},
};
use k8s_openapi::api::batch::v1::Job;
use kube::{
    Client,
    api::{Api, DeleteParams, PostParams},
    runtime::wait::{await_condition, conditions},
};
use std::sync::{Arc, atomic::AtomicUsize};
use tracing::info;
#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    let data = Arc::new(AtomicUsize::new(0));
    // build our application with a route
    let app = Router::new()
        // `GET /` goes to `root`
        .route("/", post(create_run))
        .route("/", get(get_run))
        .with_state(data.clone());

    // run our app with hyper, listening globally on port 3000
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

async fn get_run(
    // this argument tells axum to parse the request body
    // as JSON into a `CreateUser` type
    State(state): State<std::sync::Arc<AtomicUsize>>,
) -> (StatusCode, String) {
    let job_number = state.load(std::sync::atomic::Ordering::SeqCst);
    info!("Job number: {}", job_number);

    (StatusCode::CREATED, format!("Job number: {}", job_number))
}

async fn create_run(
    // this argument tells axum to parse the request body
    // as JSON into a `CreateUser` type
    State(state): State<std::sync::Arc<AtomicUsize>>,
) -> StatusCode {
    // insert your application logic here

    let Ok(client) = Client::try_default().await else {
        return StatusCode::INTERNAL_SERVER_ERROR;
    };
    let job_number = state.fetch_add(1, std::sync::atomic::Ordering::SeqCst);
    let name = format!("kir-plop-job-{}", job_number);
    let jobs: Api<Job> = Api::default_namespaced(client);

    let create_job = create_job_api_call(name.clone(), jobs.clone()).await;

    if create_job.is_err() {
        return StatusCode::INTERNAL_SERVER_ERROR;
    }

    info!("Waiting for job to complete");
    let cond = await_condition(jobs.clone(), name.as_str(), conditions::is_job_completed());
    let is_complete = tokio::time::timeout(std::time::Duration::from_secs(20), cond)
        .await
        .inspect_err(|err| {
            tracing::error!("Failed to wait for job: {:?}", err);
        });

    delete_job_api_call(name.clone(), jobs.clone())
        .await
        .expect("Failed to delete job");
    if is_complete.is_err() {
        return StatusCode::INTERNAL_SERVER_ERROR;
    }

    info!("Job completed");
    // this will be converted into a JSON response
    // with a status code of `201 Created`
    StatusCode::CREATED
}

async fn delete_job_api_call(name: String, jobs: Api<Job>) -> anyhow::Result<()> {
    info!("Delete job");

    let _delete_job = jobs.delete(&name, &DeleteParams::default()).await?;

    Ok(())
}
async fn create_job_api_call(name: String, jobs: Api<Job>) -> anyhow::Result<()> {
    info!("Creating job");
    let data = serde_json::from_value(serde_json::json!({
        "apiVersion": "batch/v1",
        "kind": "Job",
        "metadata": {
            "name": name,
        },
        "spec": {
            "template": {
                "metadata": {
                    "name": name,
                },
                "spec": {
                    "containers": [{
                        "name": "empty",
                        "image": "alpine:latest"
                    }],
                    "restartPolicy": "Never",
                }
            }
        }
    }))?;
    let _create_job = jobs.create(&PostParams::default(), &data).await?;
    Ok(())
}
