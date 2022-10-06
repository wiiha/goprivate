<script lang="ts">
   import goprivate from "./goprivate.api"
   import ed from "./encdecrypt"
   import { onMount } from "svelte"

   let notePingInfo:
      | undefined
      | { exists: boolean; consumed: boolean; consumedAt: string } = undefined
   let noteContent = ""

   /* 
   [ ] NEEDS to check if a valid key is collected
   */
   const decryptionKey = window.location.hash.substring(1)
   const noteID = window.location.pathname
      .replace("/read/", "")
      .replace(decryptionKey, "")
   console.log({ decryptionKey, noteID })

   const pingNote = async () => {
      const res = await goprivate.pingNote(noteID)
      notePingInfo = res
      console.log({ notePingInfo })
   }

   onMount(async () => {
      await pingNote()
   })

   const consumeNote = async (nid: string) => {
      const { content } = await goprivate.consumeNote(nid)
      const [iv, encryptedMessage] = content.split("$")
      const message = ed.decrypt(encryptedMessage, decryptionKey, iv)
      noteContent = message
      pingNote()
   }
</script>

<article>
   <section>
      <aside>
         {#if notePingInfo !== undefined}
            <h2>Note meta info</h2>
            <p>{`Given id: ${noteID}`}</p>
            <p>{`Exists: ${notePingInfo.exists}`}</p>
            {#if notePingInfo.exists}
               <p>{`Consumed: ${notePingInfo.consumed}`}</p>
               <p>{`Consumed at: ${notePingInfo.consumedAt}`}</p>
            {/if}
            {#if notePingInfo.exists && !notePingInfo.consumed}
               <br />
               <button on:click={() => consumeNote(noteID)}>Read note</button>
            {/if}
         {/if}
      </aside>
   </section>
   {#if noteContent !== ""}
      <section>
         <aside>
            <h2>Note</h2>
            <p>{noteContent}</p>
         </aside>
      </section>
   {/if}
</article>

<style>
</style>
