/*
 *
 * Copyright 2016 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

#include "src/core/lib/iomgr/port.h"

#ifdef GRPC_WINDOWS_SOCKETUTILS

#include "src/core/lib/iomgr/sockaddr.h"
#include "src/core/lib/iomgr/socket_utils.h"

#include <grpc/support/log.h>

const char *grpc_inet_ntop(int af, const void *src, char *dst, size_t size) {
#ifdef GPR_WIN_INET_NTOP
  return inet_ntop(af, src, dst, size);
#else
  /* Windows InetNtopA wants a mutable ip pointer */
  return InetNtopA(af, (void *)src, dst, size);
#endif /* GPR_WIN_INET_NTOP */
}

#endif /* GRPC_WINDOWS_SOCKETUTILS */