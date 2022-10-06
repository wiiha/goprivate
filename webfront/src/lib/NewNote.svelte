<script lang="ts">
   import goprivate from "./goprivate.api"

   let clearTextMsg = ""
   let urlForCreatedNote = ""

   const submitNote = async (e: SubmitEvent) => {
      console.log("got form submit", { event: e })
      if (clearTextMsg === "") return

      try {
         const res = await goprivate.newNote(clearTextMsg)
         console.log({ res })
         clearTextMsg = ""
         urlForCreatedNote = `${document.location.href}read/${res.noteID}#${res.key}`
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
         <input id="notelink" type="text" value={urlForCreatedNote} />
      </form>
   {/if}
</article>

<style>
</style>
