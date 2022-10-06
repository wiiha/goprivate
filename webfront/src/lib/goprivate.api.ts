let base = "http://localhost:8080"

if (import.meta.env.MODE === "production") {
   base = ""
}

const apiBase = `${base}/api/v1`

console.log({ apiBase })

export default {
   /* 
    Ping
  */
   ping: async function () {
      const resp = await fetch(`${apiBase}/ping`, {
         method: "GET",
         headers: {
            // Authorization: `Bearer ${getToken()}`
         }
      })
      if (resp.status !== 200) {
         throw resp
      }
      const jsonRes = await resp.json()
      return jsonRes
   },

   /* 
   New note
   */
   newNote: async function (
      clearTextContent: string
   ): Promise<{ noteID: string; }> {
      /* 
    [ ] NEEDS to encrypt note content
    */

      const resp = await fetch(`${apiBase}/newnote`, {
         method: "POST",
         headers: {
            "Content-Type": "application/json"
         },
         body: JSON.stringify({
            noteContent: clearTextContent
         })
      })
      if (resp.status !== 200) {
         throw resp
      }
      let json = await resp.json()
      return json
   },

   /* 
   Ping note to see if the given noteID
   is valid.
   */
   pingNote: async function (
      noteID: string
   ): Promise<{ exists: boolean; consumed: boolean; consumedAt: string }> {
      const resp = await fetch(`${apiBase}/pingnote/${noteID}`, {
         method: "GET",
         headers: {
            // Authorization: `Bearer ${getToken()}`
         }
      })
      if (resp.status !== 200) {
         throw resp
      }
      const jsonRes = await resp.json()
      return jsonRes
   },

   consumeNote: async function (
      noteID: string
   ): Promise<{ id: string; content: string }> {
      const resp = await fetch(`${apiBase}/consumenote/${noteID}`, {
         method: "DELETE",
         headers: {}
      })
      if (resp.status !== 200) {
         throw resp
      }
      const json = await resp.json()
      return json
   }
}
