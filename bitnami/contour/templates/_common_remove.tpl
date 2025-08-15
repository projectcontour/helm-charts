{{/*
Copyright Broadcom, Inc. All Rights Reserved.
SPDX-License-Identifier: APACHE-2.0
*/}}

{{/*
This contains helper functions from bitnami/common
However, this set of functions relates to specifically the bitnami images, and is not likely to be needed.
*/}}

{{/* vim: set filetype=mustache: */}}

{{/*
Warning about replaced images from the original.
Usage:
{{ include "common.warnings.modifiedImages" (dict "images" (list .Values.path.to.the.imageRoot) "context" $) }}
*/}}
{{- define "common.warnings.modifiedImages" -}}
{{- $affectedImages := list -}}
{{- $printMessage := false -}}
{{- $originalImages := .context.Chart.Annotations.images -}}
{{- range .images -}}
  {{- $fullImageName := printf (printf "%s/%s:%s" .registry .repository .tag) -}}
  {{- if not (contains $fullImageName $originalImages) }}
    {{- $affectedImages = append $affectedImages (printf "%s/%s:%s" .registry .repository .tag) -}}
    {{- $printMessage = true -}}
  {{- end -}}
{{- end -}}
{{- if $printMessage }}

⚠ SECURITY WARNING: Original containers have been substituted. This Helm chart was designed, tested, and validated on multiple platforms using a specific set of Bitnami and Tanzu Application Catalog containers. Substituting other containers is likely to cause degraded security and performance, broken chart features, and missing environment variables.

Substituted images detected:
{{- range $affectedImages }}
  - {{ . }}
{{- end }}
{{- end -}}
{{- end -}}


{{/*
Throw error when original container images are replaced.
The error can be bypassed by setting the "global.security.allowInsecureImages" to true. In this case,
a warning message will be shown instead.

Usage:
{{ include "common.errors.insecureImages" (dict "images" (list .Values.path.to.the.imageRoot) "context" $) }}
*/}}
{{- define "common.errors.insecureImages" -}}
{{- $relocatedImages := list -}}
{{- $replacedImages := list -}}
{{- $bitnamiLegacyImages := list -}}
{{- $retaggedImages := list -}}
{{- $globalRegistry := ((.context.Values.global).imageRegistry) -}}
{{- $originalImages := .context.Chart.Annotations.images -}}
{{- range .images -}}
  {{- $registryName := default .registry $globalRegistry -}}
  {{- $fullImageNameNoTag := printf "%s/%s" $registryName .repository -}}
  {{- $fullImageName := printf "%s:%s" $fullImageNameNoTag .tag -}}
  {{- if not (contains $fullImageNameNoTag $originalImages) -}}
    {{- if not (contains $registryName $originalImages) -}}
      {{- $relocatedImages = append $relocatedImages $fullImageName  -}}
    {{- else if not (contains .repository $originalImages) -}}
      {{- $replacedImages = append $replacedImages $fullImageName -}}
      {{- if contains "docker.io/bitnamilegacy/" $fullImageNameNoTag -}}
        {{- $bitnamiLegacyImages = append $bitnamiLegacyImages $fullImageName -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}
  {{- if not (contains (printf "%s:%s" .repository .tag) $originalImages) -}}
    {{- $retaggedImages = append $retaggedImages $fullImageName  -}}
  {{- end -}}
{{- end -}}

{{- if and (or (gt (len $relocatedImages) 0) (gt (len $replacedImages) 0)) (((.context.Values.global).security).allowInsecureImages) -}}
  {{- print "\n\n⚠ SECURITY WARNING: Verifying original container images was skipped. Please note this Helm chart was designed, tested, and validated on multiple platforms using a specific set of Bitnami and Bitnami Secure Images containers. Substituting other containers is likely to cause degraded security and performance, broken chart features, and missing environment variables.\n" -}}
{{- else if (or (gt (len $relocatedImages) 0) (gt (len $replacedImages) 0)) -}}
  {{- $errorString := "Original containers have been substituted for unrecognized ones. Deploying this chart with non-standard containers is likely to cause degraded security and performance, broken chart features, and missing environment variables." -}}
  {{- $errorString = print $errorString "\n\nUnrecognized images:" -}}
  {{- range (concat $relocatedImages $replacedImages) -}}
    {{- $errorString = print $errorString "\n  - " . -}}
  {{- end -}}
  {{- if and (eq (len $relocatedImages) 0) (eq (len $replacedImages) (len $bitnamiLegacyImages)) -}}
    {{- $errorString = print "\n\n⚠ WARNING: " $errorString -}}
    {{- print $errorString -}}
  {{- else if or (contains "docker.io/bitnami/" $originalImages) (contains "docker.io/bitnamiprem/" $originalImages) (contains "docker.io/bitnamisecure/" $originalImages) -}}
    {{- $errorString = print "\n\n⚠ ERROR: " $errorString -}}
    {{- $errorString = print $errorString "\n\nIf you are sure you want to proceed with non-standard containers, you can skip container image verification by setting the global parameter 'global.security.allowInsecureImages' to true." -}}
    {{- $errorString = print $errorString "\nFurther information can be obtained at https://github.com/bitnami/charts/issues/30850" -}}
    {{- print $errorString | fail -}}
  {{- else if gt (len $replacedImages) 0 -}}
    {{- $errorString = print "\n\n⚠ WARNING: " $errorString -}}
    {{- print $errorString -}}
  {{- end -}}
{{- else if gt (len $retaggedImages) 0 -}}
  {{- $warnString := "\n\n⚠ WARNING: Original containers have been retagged. Please note this Helm chart was tested, and validated on multiple platforms using a specific set of Bitnami and Bitnami Secure Images containers. Substituting original image tags could cause unexpected behavior." -}}
  {{- $warnString = print $warnString "\n\nRetagged images:" -}}
  {{- range $retaggedImages -}}
    {{- $warnString = print $warnString "\n  - " . -}}
  {{- end -}}
  {{- print $warnString -}}
{{- end -}}
{{- end -}}
