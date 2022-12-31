# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM us-central1-docker.pkg.dev/cloud-workstations-images/predefined/code-oss:latest

RUN curl -L -o /usr/bin/bazel https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-amd64 && chmod a+x /usr/bin/bazel
RUN curl -L -o /usr/bin/buildifier https://github.com/bazelbuild/buildtools/releases/latest/download/buildifier-linux-amd64 && chmod a+x /usr/bin/buildifier
RUN curl -L -o /usr/bin/buildozer https://github.com/bazelbuild/buildtools/releases/latest/download/buildozer-linux-amd64 && chmod a+x /usr/bin/buildozer
RUN curl -L -o /usr/bin/unused_deps https://github.com/bazelbuild/buildtools/releases/latest/download/unused_deps-linux-amd64 && chmod a+x /usr/bin/unused_deps

RUN wget https://open-vsx.org/api/BazelBuild/vscode-bazel/0.7.0/file/BazelBuild.vscode-bazel-0.7.0.vsix && unzip BazelBuild.vscode-bazel-0.7.0.vsix "extension/*" && mv extension /opt/code-oss/extensions/bazel
