package ch.bfh.ti.hirtp1ganzg1.thesis

import ch.bfh.ti.hirtp1ganzg1.thesis.api.services.impl.ICacheDefaultImpl
import org.junit.Test
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertTrue

class CacheTest {
    @Test
    fun testCache() {
        val cache =
            ICacheDefaultImpl<String, String>()
        val key = "key"
        val value = "value"
        cache.set(key, value)
        assertTrue(cache.exists(key), key)
        assertEquals(value, cache.get(key))
        cache.remove(key)
        assertFalse(cache.exists(key), key)
    }
}
