<script lang="ts">
   import goprivate from "./goprivate.api"
   import ed from "./encdecrypt"

   let clearTextMsg = ""
   let urlForCreatedNote = ""
   let linkCopiedToClipboard = false

   const submitNote = async (e: SubmitEvent) => {
      console.log("got form submit", { event: e })
      if (clearTextMsg === "") return

      const { encryptedHex, ivHex, keyHex } = await ed.encrypt(clearTextMsg)
      try {
         const res = await goprivate.newNote(`${ivHex}\$${encryptedHex}`)
         console.log({ res })
         clearTextMsg = ""
         urlForCreatedNote = `${document.location.href}read/${res.noteID}#${keyHex}`
      } catch (error) {
         console.log({ error })
      }
   }
</script>

<article>
   <h1>New Note</h1>
   <form on:submit|preventDefault={submitNote}>
      <textarea
         bind:value={clearTextMsg}
         style="resize: vertical;"
         cols="30"
         rows="10"
      />
      <button on:click|stopPropagation type="submit">Create note</button>
   </form>
   <br />
   {#if urlForCreatedNote !== ""}
      <form on:submit|preventDefault={() => {}}>
         <label for="notelink"><h2>Note created!</h2></label>
         <p>Your note can be accessed using this link.</p>
         <input
            on:click={(ev) => {
               const link = ev.currentTarget.value
               console.log({ link })
               navigator.clipboard.writeText(link)
               linkCopiedToClipboard = true
            }}
            id="notelink"
            type="text"
            value={urlForCreatedNote}
         />
         {#if linkCopiedToClipboard}
         <p>Link is copied to clipboard!</p>
         {/if}
      </form>
   {/if}
</article>

<style>
</style>
