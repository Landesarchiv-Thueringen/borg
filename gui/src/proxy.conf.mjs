import * as dotenv from "dotenv";

dotenv.config({ path: "../.env" });

export default {
  "/analyze-file": {
    target: `http://localhost:${process.env.PORT}`,
    secure: false,
  },
};
