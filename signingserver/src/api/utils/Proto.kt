package ch.bfh.ti.hirtp1ganzg1.thesis.api.utils

import com.google.protobuf.ByteString

fun ByteArray.toByteString(): ByteString = ByteString.copyFrom(this)