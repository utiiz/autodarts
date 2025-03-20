// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
// Navigation component

package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Navbar() templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!-- ========== HEADER ========== --><header class=\"sticky top-0 inset-x-0 flex flex-wrap md:justify-start md:flex-nowrap z-50 w-full text-sm\"><nav class=\"mt-4 relative max-w-2xl w-full bg-white border border-gray-200 rounded-full mx-2 py-2.5 md:flex md:items-center md:justify-between md:py-0 md:px-4 md:mx-auto\"><div class=\"px-4 md:px-0 flex justify-between items-center w-full py-2\"><div><a href=\"/\" @click.prevent=\"loadContent(&#39;/&#39;)\" :class=\"route === &#39;/&#39; ? &#39;font-bold text-gray-800&#39; : &#39;text-gray-500&#39;\" class=\"flex gap-x-2 items-center text-base inline-block focus:outline-hidden focus:opacity-80\">GoCards</a></div><div><a href=\"/game\" @click.prevent=\"loadContent(&#39;/game&#39;)\" :class=\"route === &#39;/game&#39; ? &#39;font-bold text-gray-800&#39; : &#39;text-gray-500&#39;\" class=\"flex gap-x-2 items-center text-base inline-block focus:outline-hidden focus:opacity-80\">Game</a></div><div><a href=\"/dashboard\" @click.prevent=\"loadContent(&#39;/dashboard&#39;)\" :class=\"route === &#39;/dashboard&#39; ? &#39;font-bold text-gray-800&#39;\n\t\t\t\t\t\t\t: &#39;text-gray-500&#39;\" class=\"py-0.5 md:py-3 px-4 md:px-1 text-gray-500 hover:text-gray-800 focus:outline-hidden\" aria-current=\"page\">Dashboard</a></div></div></nav></header><!-- ========== END HEADER ========== -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
