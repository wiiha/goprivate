<script lang="ts">
   import goprivate from "./goprivate.api"
   import ed from "./encdecrypt"
   import { onMount } from "svelte"

   let notePingInfo:
      | undefined
      | { exists: boolean; consumed: boolean; consumedAt: string } = undefined
   let noteContent = ""
   let noteContentIsValid = true
   let missingKeyInUrl = false
   let userDefinedPassword = ""

   /* 
   [ ] NEEDS to check if a valid key is collected
   */
   const decryptionKey = window.location.hash.substring(1)
   missingKeyInUrl = decryptionKey === ""
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
      if (missingKeyInUrl && userDefinedPassword === "") return
      const { content } = await goprivate.consumeNote(nid)
      const [iv, encryptedMessage] = content.split("$")
      const res = missingKeyInUrl
         ? await ed.decryptWithPassword(
              encryptedMessage,
              userDefinedPassword,
              iv
           )
         : await ed.decrypt(encryptedMessage, decryptionKey, iv)
      noteContent = res.message
      noteContentIsValid = res.messageIsValid
      missingKeyInUrl = false
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
            {#if !notePingInfo.exists}
               <p><b>This note does not exist.</b></p>
            {/if}
            {#if notePingInfo.exists && !notePingInfo.consumed && missingKeyInUrl}
               <p><b>Important</b></p>
               <p>
                  Please enter the password. The author should have given it to
                  you. <b>The note is lost if you enter the wrong password.</b>
               </p>
               <input
                  bind:value={userDefinedPassword}
                  type="text"
                  name="userPassword"
                  id="userPassword"
               />
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
            {#if !noteContentIsValid}
               <p>
                  <mark
                     >Clear text message could not be validated. You might have
                     entered the wrong key. This cannot be reversed, the message
                     is lost.
                  </mark>
               </p>
               <p>Decrypted but invalid message:</p>
            {/if}
            <p>{noteContent}</p>
         </aside>
      </section>
   {/if}
</article>

<style>
</style>
