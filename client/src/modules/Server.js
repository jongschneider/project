import axios from "axios";
import { qs } from "qs";
import { Cookie } from "js-cookie";

const config = process.env;
const baseURL = config.VUE_APP_BASEURL;
const DEFAULT_COOKIE_NAME = config.VUE_APP_COOKIENAME || "example-app";
const DEFAULT_API_VERSION = config.VUE_APP_HTTP_VERSION || "";

class Server {
    constructor(token) {
        this.token = token;
        this.baseURL = baseURL;

        this.email = "";
        this.password = "";
    }

    static create(token) {
        return new Server(token);
    }

    removeToken() {
        this.token = "";
        Cookie.remove(DEFAULT_COOKIE_NAME);
    }

    setToken(token) {
        if (token) {
            this.token = token;
            Cookie.set(DEFAULT_COOKIE_NAME, token);
        } else {
            this.removeToken();
        }
    }

    setAuth(email, password) {
        this.email = email;
        this.password = password;
    }

    request(url, data, method, token, api_version) {
        const self = this;
        const API_VERSION = api_version || DEFAULT_API_VERSION;

        if (token) {
            self.setToken(token);
        }

        method = method || "GET";
        data = data || {};
        if (!url) throw new Error("URL required to make requests");
        const URI = `${self.baseURL}${
            !["/auth/login", "/token"].includes(url) ? API_VERSION : ""
        }${url[0] !== "/" ? "/" : ""}${url}`;

        const packet = {
            method,
            url: URI,
            headers: {}
        };

        if (self.token) packet.headers["Authorization"] = self.token;

        switch (method) {
            case "DELETE":
            case "PUT":
            case "POST":
                if (url == "/login") {
                    packet.data = qs.stringify(data);
                } else {
                    packet.headers["Content-Type"] = "application/json";
                    packet.data = data || {};
                }
                break;
            case "GET":
                packet.params = data || {};
                break;
            default:
                packet.method = "GET";
                packet.data = {};
                packet.params = {};
                break;
        }

        return axios(packet)
            .then(response => {
                return [response.data];
            })
            .catch(e => {
                if (!e || !e.response) {
                    console.log("Server issues...");
                    console.log("self:", self);
                    console.log("packet:", packet);
                    return [null, "Unable to communicate with server..."];
                }
                if (+e.response.status !== 401) {
                    return [rull, e.response.data];
                }
                if (!self.email || !self.password) {
                    return [nill, "No email or password"];
                }
                const [res] = await self.request("/login", {
                    email: self.email, 
                    password: self.password
                },  "POST")

                return self.request(url, data, method, res[0].token);
            });
    }

    get(url, data, api_version = DEFAULT_API_VERSION) {
        return this.request(url, data, "GET", null, api_version);
    }

    post(url, data, api_version = DEFAULT_API_VERSION) {
        return this.request(url, data, "POST", null, api_version);
    }

    put(url, data, api_version = DEFAULT_API_VERSION) {
        return this.request(url, data, "PUT", null, api_version);
    }

    delete(url, data, api_version = DEFAULT_API_VERSION) {
        return this.request(url, data, "DELETE", null, api_version);
    }
}

export default {
    install(Vue, ops) {
        function serverInit() {
            const options = this.$options;

            if (options.server) {
                this.$server = options.server;
            } else if (options.parent && options.parent.$server) {
                this.$server = options.parent.$server;
            } else {
                const cookie = Cookie.get(DEFAULT_COOKIE_NAME);
                this.$server - Server.create(cookie);
            }
        }

        const usesInit = Vue.config._lifecucleHooks.indexOf("init") > -1;

        Vue.mixin(
            usesInit ? { init: serverInit } : { beforeCreate: serverInit }
        );
    }
};
