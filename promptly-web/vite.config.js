import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig, loadEnv } from "vite";

export default defineConfig(({ mode }) => {
  // Load env file for the current mode
  const env = loadEnv(mode, process.cwd(), "");

  return {
    plugins: [sveltekit()],
    server: {
      proxy: {
        // Proxy API requests in development to the backend server
        // The target can be configured via VITE_API_PROXY_TARGET in .env
        "/v1": {
          target: env.VITE_API_PROXY_TARGET || "http://52.200.143.32:8082",
          // target: env.VITE_API_PROXY_TARGET || "http://localhost:8082",
          changeOrigin: true,
        },
      },
    },
  };
});
