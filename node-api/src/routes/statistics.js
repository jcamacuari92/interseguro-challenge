const express = require('express');
const { validateMatrices, computeStatistics } = require('../services/statisticsService');

const router = express.Router();

router.post('/statistics', (req, res) => {
  try {
    const { matrices } = req.body || {};
    validateMatrices(matrices);
    const stats = computeStatistics(matrices);
    res.json(stats);
  } catch (err) {
    res.status(400).json({ error: err.message });
  }
});

module.exports = router;