// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "server/views/components"

func DashboardPage() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div x-data=\"{\n\t\t\tauth: pb.authStore,\n\t\t\tloading: false,\n\t\t\tasync logout() {\n\t\t\t\tthis.loading = true;\n\t\t\t\ttry {\n\t\t\t\t\tawait pb.authStore.clear();\n\t\t\t\t\tloadContent(&#39;/login&#39;);\n\t\t\t\t} catch (error) {\n\t\t\t\t\tconsole.log(error);\n\t\t\t\t} finally {\n\t\t\t\t\tthis.loading = false;\n\t\t\t\t}\n\t\t\t}\n\t\t}\" x-init=\"if (!auth.isValid) loadContent(&#39;/login&#39;);\"><template x-if=\"auth.isValid\"><div><!-- Hero --><div class=\"relative overflow-hidden before:absolute before:top-0 before:start-1/2 before:bg-[url(&#39;https://preline.co/assets/svg/examples/squared-bg-element.svg&#39;)] before:bg-no-repeat before:bg-top before:size-full before:-z-10 before:transform before:-translate-x-1/2\"><div class=\"max-w-[85rem] mx-auto px-4 sm:px-6 lg:px-8 pt-24 pb-10\"><!-- Title --><div class=\"mt-5 max-w-xl text-center mx-auto\"><h1 class=\"block font-bold text-gray-800 text-4xl md:text-5xl lg:text-6xl\">Welcome to the Dashboard</h1></div><!-- End Title --><div class=\"mt-5 max-w-3xl text-center mx-auto\"><p class=\"text-lg text-gray-600\" x-text=\"auth.model.email\"></p></div><!-- Buttons --><div class=\"mt-8 gap-3 flex justify-center\"><a class=\"\n\t\t\t\t\t\t\t\t\tinline-flex justify-center items-center gap-x-3 \n\t\t\t\t\t\t\t\t\ttext-center text-white text-sm font-medium\n\t\t\t\t\t\t\t\t\tbg-gradient-to-tl from-blue-600 to-violet-600\n\t\t\t\t\t\t\t\t\thover:from-violet-600 hover:to-blue-600\n\t\t\t\t\t\t\t\t\tfocus:from-violet-600 focus:to-blue-600\n\t\t\t\t\t\t\t\t\tborder border-transparent rounded-full py-3 px-4\" href=\"/logout\" @click.prevent=\"logout\">Log Out</a></div><!-- End Buttons --></div></div><!-- End Hero --><!-- Card Blog --><div class=\"max-w-[85rem] px-4 py-10 sm:px-6 lg:px-8 lg:py-14 mx-auto\"><!-- Title --><div class=\"max-w-2xl mx-auto text-center mb-10 lg:mb-14\"><h2 class=\"text-2xl font-bold md:text-4xl md:leading-tight\">My Cards</h2><p class=\"mt-1 text-gray-600\">Latest</p></div><!-- End Title --><!-- Grid --><div x-data=\"{\n\t\t\t\t\t\tloading: false,\n\t\t\t\t\t\titems: [],\n\t\t\t\t\t\tpage: 0,\n\t\t\t\t\t\ttotalPages: 1,\n\t\t\t\t\t\tasync loadItems() {\n\t\t\t\t\t\t\tthis.loading=true;\n\t\t\t\t\t\t\ttry {\n\t\t\t\t\t\t\t\tthis.page++;\n\t\t\t\t\t\t\t\tif (this.page &gt; this.totalPages) return;\n\n\t\t\t\t\t\t\t\tconst response = await pb.collection(&#39;cards&#39;).getList(this.page, 9, {\n\t\t\t\t\t\t\t\t\tsort: &#39;-created&#39;,\n\t\t\t\t\t\t\t\t\tfilter: `user = &#39;${pb.authStore.model.id}&#39;`,\n\t\t\t\t\t\t\t\t\texpand: &#39;user&#39;\n\t\t\t\t\t\t\t\t});\n\n\t\t\t\t\t\t\t\tthis.items = [...this.items, ...response.items];\n\t\t\t\t\t\t\t\tthis.totalPages = response.totalPages\n\t\t\t\t\t\t\t} catch (error) {\n\t\t\t\t\t\t\t\tconsole.log(error);\n\t\t\t\t\t\t\t} finally {\n\t\t\t\t\t\t\t\tthis.loading = false;\n\t\t\t\t\t\t\t}\n\t\t\t\t\t\t}\n\t\t\t\t\t}\" x-init=\"$nextTick(() =&gt; loadItems())\"><div class=\"grid sm:grid-cols-2 lg:grid-cols-3 gap-6\"><template x-for=\"item in items\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.Card().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</template></div><!-- End Grid --><!-- Card --><div class=\"mt-12 text-center\" :disabled=\"loading\" x-show=\"page &lt; totalPages\"><button @click=\"loadItems\" class=\"py-3 px-4 inline-flex items-center gap-x-1 text-sm font-medium rounded-full border border-gray-200 bg-white text-blue-600 shadow-2xs hover:bg-gray-50 focus:outline-hidden focus:bg-gray-50 disabled:opacity-50 disabled:pointer-events-none\">Load more</button></div><!-- End Card --></div><!-- End Card Blog --></div></div></template></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
