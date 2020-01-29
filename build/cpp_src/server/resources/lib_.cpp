	#include <cmrc/cmrc.hpp>
	#include <map>

	namespace cmrc { namespace resources {

	namespace res_chars {
	// These are the files which are available in this resource library
	
	}

	inline void load_resources() {
	    // This initializes the list of resources and pointers to their data
	    static std::once_flag flag;
	    std::call_once(flag, [] {
		cmrc::detail::with_table([](resource_table& table) {
			(void)table;
		    
		});
	    });
	}

	// namespace {
	//     extern struct resource_initializer {
	//         resource_initializer() {
	//             load_resources();
	//         }
	//     } dummy;
	// }

	}}

	// The resource library initialization function. Intended to be called
	// before anyone intends to use any of the resource defined by this
	// resource library
	extern void cmrc_init_resources_resources() {
	    cmrc::resources::load_resources();
	}
    