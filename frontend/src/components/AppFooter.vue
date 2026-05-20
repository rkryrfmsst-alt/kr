<template>
  <footer class="footer">
    <div class="container">
      <span class="copyright">© {{ year }} Solutions</span>

      <div class="socials">
        <a href="https://t.me/" target="_blank" rel="noopener" class="social-link">
          <img src="../assets/icons/telegram.svg" alt="Telegram" class="social-icon">
        </a>
      </div>
    </div>
  </footer>
</template>

<script setup>
const year = new Date().getFullYear();
</script>

<style scoped>
.footer {
  background: #000;
  font-family: 'Inter', sans-serif;
  width: 100%;
  padding: 20px 0;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 50px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.copyright {
  color: rgba(255, 255, 255, 0.4);
  font-size: 0.85rem;
}

.social-link {
  display: flex;
  align-items: center;
  opacity: 0.5;
  transition: opacity 0.2s;
}
.social-link:hover { opacity: 1; }

.social-icon {
  width: 20px;
  height: 20px;
}
</style>