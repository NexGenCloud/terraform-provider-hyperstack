variable "name" {
  type = string
}

variable "namespace" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "artifacts_dir" {
  type = string
}

variable "kube_config_file" {
  type = string
}

variable "load_balancer_address" {
  type = string
}

variable "vllm_model_name" {
  description = "The model name to be used in the deployment"
  type        = string
  default     = "NousResearch/Meta-Llama-3-8B-Instruct"
}
