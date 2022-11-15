from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

APIRatelimit: RatelimitType
DESCRIPTOR: _descriptor.FileDescriptor
IPRatelimit: RatelimitType
InstanceRatelimit: RatelimitType
ServiceRatelimit: RatelimitType

class RateLimitPluginRequest(_message.Message):
    __slots__ = ["key", "type"]
    KEY_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    key: str
    type: RatelimitType
    def __init__(self, type: _Optional[_Union[RatelimitType, str]] = ..., key: _Optional[str] = ...) -> None: ...

class RateLimitPluginResponse(_message.Message):
    __slots__ = ["allow"]
    ALLOW_FIELD_NUMBER: _ClassVar[int]
    allow: bool
    def __init__(self, allow: bool = ...) -> None: ...

class RatelimitType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = []
