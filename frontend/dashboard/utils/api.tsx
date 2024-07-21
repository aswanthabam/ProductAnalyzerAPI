export function Get(
  input: string | URL | globalThis.Request,
  init?: RequestInit
) {
  if (init == null) {
    init = {};
  }
  if (init?.headers == null) {
    init = { ...init, headers: {} };
  }
  init!.headers = {
    ...init!.headers,
    Authorization: `Bearer ${localStorage.getItem("token")}`,
  };
  init!.method = "GET";
  return fetch(input, init);
}

export function Post(
  input: string | URL | globalThis.Request,
  init?: RequestInit
) {
  if (init == null) {
    init = {};
  }
  if (init?.headers == null) {
    init = { ...init, headers: {} };
  }
  init!.headers = {
    ...init!.headers,
    Authorization: `Bearer ${localStorage.getItem("token")}`,
  };
  init!.method = "POST";
  return fetch(input, init);
}

export function IsAuthenticated() {
  var res = localStorage.getItem("token") != null;
  console.log(res);
  return res;
}
