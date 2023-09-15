// literals
export const QUERY_KEY = "rpcUrl";
export const costPerByte = 1;

export const proxyUrl = import.meta.env.VITE_PROXY_URL;
export const fileRelativePath = "People/mimasa/test/imgformat/img/w3c_home.png";
export const fileUrl = proxyUrl + fileRelativePath;
export const dataSize = 6833;

// env vars
export const provider = import.meta.env.VITE_PROVIDER;
export const hub = import.meta.env.VITE_HUB;
export const defaultNitroRPCUrl = import.meta.env.VITE_NITRO_RPC_URL;
export const initialChannelBalance = parseInt(
  import.meta.env.VITE_INITIAL_CHANNEL_BALANCE,
  10
);
