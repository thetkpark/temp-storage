const express = require("express");
const app = express();
const cors = require("cors");
const helmet = require("helmet");
app.use(
  cors({
    origin: ["http://localhost:5500",'*'],
    methods: ["GET", "PUT", "POST", "PATCH", "DELETE"],
  })
);
app.use(helmet());
app.use(express.static("client"));
const PORT = process.env.PORT || 5000;
app.post("/upload", (req, res) => {
  res.json({ url: "https://aka.cscms.me/test" });
});
app.get("/:shortUrl", (req, res) => {
  // For send the file back to user
});
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
