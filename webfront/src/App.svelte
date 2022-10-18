<script lang="ts">
   import NewNote from "./lib/NewNote.svelte"
   import ConsumeNote from "./lib/ConsumeNote.svelte"

   console.log("MODE: ", import.meta.env.MODE)

   let aboutSectionIsOpen = true

   // Mini router setup
   const pageLanding = "landing"
   const pageConsume = "consume"
   let currentPage = pageLanding
   // end router setup

   let urlPath = window.location.pathname

   if (/\/read\/\w+/.test(urlPath))
      (currentPage = pageConsume) && (aboutSectionIsOpen = false)
   if (/^\/read\/$/.test(urlPath)) window.location.replace("/")
</script>

<header>
   <h1>
      <b>
         <a id="siteLogo" href="/"> GoPrivate </a>
      </b>
   </h1>
</header>

<hr class="m-0" />

<main>
   <article>
      <aside>
         <details id="about" bind:open={aboutSectionIsOpen}>
            <summary>About this service</summary>
            <p>
               <b>GoPrivate</b> allows you to write a message that will be encrypted
               before it is sent and stored on the server.
            </p>
            <p>
               A link containing the password will be generated in the browser.
               Share this link with someone in order for them to read the
               message.
            </p>
            <p>
               The message can only be read once and its content will thereafter
               be deleted from the server.
            </p>
            <p>
               However, a record of the message ID and when it was read will be
               stored. This way you can see if someone has opened your message.
            </p>
         </details>
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
      Created by <a
         href="https://github.com/wiiha"
         target="_blank"
         rel="noopener noreferrer">wiiha</a
      >😇
   </p>
</footer>

<style>
   #siteLogo {
      color: black;
   }

   #about > summary {
      font-size: 1.2rem;
   }
</style>
