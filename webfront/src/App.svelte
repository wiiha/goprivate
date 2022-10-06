<script lang="ts">
   import NewNote from "./lib/NewNote.svelte"
   import ConsumeNote from "./lib/ConsumeNote.svelte"

   console.log("MODE: ", import.meta.env.MODE)

   // Mini router setup
   const pageLanding = "landing"
//    const pagePing = "ping"
   const pageConsume = "consume"
   let currentPage = pageLanding
   // end router setup

   let urlPath = window.location.pathname

   if (urlPath.startsWith("/read/")) currentPage = pageConsume
</script>

<main>
   {#if currentPage == pageLanding}
      <NewNote />
   {:else if currentPage == pageConsume}
      <ConsumeNote />
   {/if}
</main>

<style>
</style>
