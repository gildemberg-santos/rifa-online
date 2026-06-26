<script setup lang="ts">
import Navbar from "./layouts/Navbar.vue"
import Footer from "./layouts/Footer.vue"
import NotificationBar from "./components/NotificationBar.vue"
</script>

<template>
  <div class="flex flex-col min-h-screen">
    <NotificationBar />
    <Navbar />
    <main class="flex-1">
      <router-view />
    </main>
    <Footer />
  </div>
</template>
