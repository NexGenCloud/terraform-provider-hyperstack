FROM golang:1.23 as runtime

ENV TZ=UTC
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt update \
 && apt install -y --no-install-recommends \
      jq \
      yq

RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go install github.com/goreleaser/goreleaser/v2@latest
RUN go install github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi@latest
RUN go install github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework@latest

COPY build/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]