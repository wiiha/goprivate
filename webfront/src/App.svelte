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

<header>
   <h1>
      <b> GoPrivate </b>
   </h1>
</header>

<hr class="m-0" />

<main>
   <article>
      <aside>
         <p>This is a short text describing the service.</p>
      </aside>
   </article>

   {#if currentPage == pageLanding}
      <NewNote />
   {:else if currentPage == pageConsume}
      <ConsumeNote />
   {/if}
</main>

<hr class="m-0" />

<footer>
   <p>
      <i>
         <mark>This is the footer area</mark>
      </i>
   </p>
</footer>

<style>
</style>
