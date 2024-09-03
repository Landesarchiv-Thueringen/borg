import * as dotenv from "dotenv";

dotenv.config({ path: "../.env" });

export default {
  "/api": {
    target: `http://localhost:${process.env.PORT}`,
    secure: false,
  },
};
