{{- /*
Common template helpers
*/ -}}

{{- define "contour-gateway-provisioner.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "contour-gateway-provisioner.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := include "contour-gateway-provisioner.name" . -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "contour-gateway-provisioner.labels" -}}
app.kubernetes.io/name: {{ include "contour-gateway-provisioner.name" . }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: gateway-provisioner
{{- end -}}

{{- define "contour-gateway-provisioner.selectorLabels" -}}
app.kubernetes.io/name: {{ include "contour-gateway-provisioner.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "contour-gateway-provisioner.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (include "contour-gateway-provisioner.fullname" .) .Values.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{- define "contour-gateway-provisioner.serviceAccountNamespace" -}}
{{- default .Release.Namespace .Values.serviceAccount.namespace -}}
{{- end -}}
