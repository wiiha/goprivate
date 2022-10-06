<script lang="ts">
   import NewNote from "./lib/NewNote.svelte"
   import goprivate from "./lib/goprivate.api"

   console.log("MODE: ", import.meta.env.MODE)

   // Mini router setup
   const pageLanding = "landing"
   const pagePing = "ping"
   const pageConsume = "consume"
   let currentPage = pageLanding
   // end router setup

   let urlPath = window.location.pathname

   if (urlPath.startsWith("/read/")) console.log("yeap!")
</script>

<main>
   {#if currentPage == pageLanding}
      <NewNote />
   {/if}
</main>

<style>
</style>
